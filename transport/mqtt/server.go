package mqtt

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bytectlgo/protoc-gen-gomq/pkg/matcher"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	xhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

// ServerOption is an MQTT server option.
type ServerOption func(*Server)

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.middleware.Use(m...)
	}
}

// Filter with MQTT middleware option.
func Filter(filters ...xhttp.FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

// RequestVarsDecoder with request decoder.
func RequestVarsDecoder(dec xhttp.DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decVars = dec
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec xhttp.DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decBody = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en xhttp.EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = en
	}
}

// ErrorEncoder with error encoder.
func ErrorEncoder(en xhttp.EncodeErrorFunc) ServerOption {
	return func(o *Server) {
		o.ene = en
	}
}

// StrictSlash is with mux's StrictSlash
// If true, when the path pattern is "/path/", accessing "/path" will
// redirect to the former and vice versa.
func StrictSlash(strictSlash bool) ServerOption {
	return func(o *Server) {
		o.strictSlash = strictSlash
	}
}

func NotFoundHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.NotFoundHandler = handler
	}
}

func MethodNotAllowedHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.MethodNotAllowedHandler = handler
	}
}

func WithClientOption(option *mqtt.ClientOptions) ServerOption {
	return func(s *Server) {
		s.clientOption = option
	}
}

// Server is an MQTT  Topic Route server wrapper.
type Server struct {
	err          error
	timeout      time.Duration
	filters      []xhttp.FilterFunc
	middleware   matcher.Matcher
	decVars      xhttp.DecodeRequestFunc
	decBody      xhttp.DecodeRequestFunc
	enc          xhttp.EncodeResponseFunc
	ene          xhttp.EncodeErrorFunc
	strictSlash  bool
	router       *mux.Router
	Handler      http.Handler
	clientOption *mqtt.ClientOptions
	mqttClient   mqtt.Client
}

// NewServer creates an MQTT server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		clientOption: mqtt.NewClientOptions(),
		timeout:      1 * time.Second,
		middleware:   matcher.New(),
		decVars:      xhttp.DefaultRequestVars,
		decBody:      xhttp.DefaultRequestDecoder,
		enc:          xhttp.DefaultResponseEncoder,
		ene:          xhttp.DefaultErrorEncoder,
		strictSlash:  true,
		router:       mux.NewRouter(),
	}

	srv.router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	for _, o := range opts {
		o(srv)
	}
	srv.router.StrictSlash(srv.strictSlash)
	srv.router.Use(srv.filter())
	srv.Handler = xhttp.FilterChain(srv.filters...)(srv.router)
	srv.mqttClient = mqtt.NewClient(srv.clientOption)
	return srv
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Warnf("mqtt route not found: %s", r.URL.Path)
}
func (s *Server) MQTTClient() mqtt.Client {
	return s.mqttClient
}
func (s *Server) MQTTHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		rw := &MQTTResponseWriter{
			header: http.Header{},
			client: client,
		}
		req := &http.Request{
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Method: http.MethodPost,
			URL:    &url.URL{Path: getUrlPathFromTopic(msg.Topic())},
			Body:   io.NopCloser(bytes.NewReader(msg.Payload())),
		}
		ctx := WithClient(req.Context(), client)
		ctx = WithMessage(ctx, msg)
		req = req.WithContext(ctx)
		s.Handler.ServeHTTP(rw, req)
	}
}

func getUrlPathFromTopic(topic string) string {
	if strings.HasPrefix(topic, "/") {
		return topic
	}
	return "/noSlashRoot/" + topic
}

// Use uses a service middleware with selector.
// selector:
//   - '/*'
//   - '/helloworld.v1.Greeter/*'
//   - '/helloworld.v1.Greeter/SayHello'
func (s *Server) Use(selector string, m ...middleware.Middleware) {
	s.middleware.Add(selector, m...)
}

// WalkRoute walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkRoute(fn WalkRouteFunc) error {
	return s.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return nil // ignore no methods
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		for _, method := range methods {
			if err := fn(RouteInfo{Method: method, Path: path}); err != nil {
				return err
			}
		}
		return nil
	})
}

// WalkHandle walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkHandle(handle func(method, path string, handler http.HandlerFunc)) error {
	return s.WalkRoute(func(r RouteInfo) error {
		handle(r.Method, r.Path, s.ServeHTTP)
		return nil
	})
}

// Route registers an MQTT router.
func (s *Server) Route(prefix string, filters ...xhttp.FilterFunc) *Router {
	prefix = getUrlPathFromTopic(prefix)
	return newRouter(prefix, s, filters...)
}

// ServeMQTT should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}

func (s *Server) filter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(req.Context(), s.timeout)
			} else {
				ctx, cancel = context.WithCancel(req.Context())
			}
			defer cancel()

			pathTemplate := req.URL.Path
			if route := mux.CurrentRoute(req); route != nil {
				// /path/123 -> /path/{id}
				pathTemplate, _ = route.GetPathTemplate()
			}
			msg, _ := MessageFromServerContext(req.Context())
			client, _ := ClientFromServerContext(req.Context())
			tr := &Transport{
				operation:   pathTemplate,
				reqHeader:   headerCarrier(req.Header),
				replyHeader: headerCarrier(w.Header()),
				request:     req,
				response:    w,
				message:     msg,
				client:      client,
			}
			tr.request = req.WithContext(transport.NewServerContext(ctx, tr))
			next.ServeHTTP(w, tr.request)
		})
	}
}

// Start start the MQTT server.
func (s *Server) Start(ctx context.Context) error {
	log.Info("[MQTT] server starting")
	token := s.mqttClient.Connect()
	if !token.WaitTimeout(s.timeout) {
		log.Errorf("mqtt connect wait timeout, address: %s", s.clientOption.Servers[0].String())
	}
	if token.Error() != nil {
		log.Errorf("mqtt connect error, %s", token.Error().Error())
		panic(token.Error())
	}
	return nil
}

// Stop stop the MQTT server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info("[MQTT] server stopping")
	if s.mqttClient != nil {
		s.mqttClient.Disconnect(1000)
	}
	return s.Shutdown(ctx)
}

func (s *Server) Endpoint() (*url.URL, error) {
	return &url.URL{Scheme: "mqtt", Host: s.clientOption.Servers[0].String()}, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}
