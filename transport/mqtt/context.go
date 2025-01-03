package mqtt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ Context = (*wrapper)(nil)

// Context is an MQTT Context.
type Context interface {
	context.Context
	Vars() url.Values
	Client() mqtt.Client
	Message() mqtt.Message
	Response() http.ResponseWriter
	Middleware(middleware.Handler) middleware.Handler
	Bind(interface{}) error
	BindVars(interface{}) error
	JSON(interface{}) error
	String(string) error
	Stream(string, io.Reader) error
	Reset(http.ResponseWriter, *http.Request)
}

type responseWriter struct {
	code int
	w    http.ResponseWriter
}

func (w *responseWriter) reset(res http.ResponseWriter) {
	w.w = res
	w.code = http.StatusOK
}
func (w *responseWriter) Header() http.Header        { return w.w.Header() }
func (w *responseWriter) WriteHeader(statusCode int) { w.code = statusCode }
func (w *responseWriter) Write(data []byte) (int, error) {
	w.w.WriteHeader(w.code)
	return w.w.Write(data)
}
func (w *responseWriter) Unwrap() http.ResponseWriter { return w.w }

type wrapper struct {
	router *Router
	req    *http.Request
	res    http.ResponseWriter
	w      responseWriter
}

func (c *wrapper) Header() http.Header {
	return c.req.Header
}

func (c *wrapper) Vars() url.Values {
	raws := mux.Vars(c.req)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return vars
}

func (c *wrapper) Form() url.Values {
	if err := c.req.ParseForm(); err != nil {
		return url.Values{}
	}
	return c.req.Form
}

func (c *wrapper) Query() url.Values {
	return c.req.URL.Query()
}

func (c *wrapper) Message() mqtt.Message { return MessageFromContext(c.req.Context()) }

func (c *wrapper) Response() http.ResponseWriter { return c.res }

func (c *wrapper) Middleware(h middleware.Handler) middleware.Handler {
	if tr, ok := transport.FromServerContext(c.req.Context()); ok {
		return middleware.Chain(c.router.srv.middleware.Match(tr.Operation())...)(h)
	}
	return middleware.Chain(c.router.srv.middleware.Match(c.req.URL.Path)...)(h)
}
func (c *wrapper) Bind(v interface{}) error     { return c.router.srv.decBody(c.req, v) }
func (c *wrapper) BindVars(v interface{}) error { return c.router.srv.decVars(c.req, v) }

func (c *wrapper) Returns(v interface{}, err error) error {
	if err != nil {
		return err
	}
	return c.router.srv.enc(&c.w, c.req, v)
}

func (c *wrapper) JSON(v interface{}) error {
	c.res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.res).Encode(v)
}

func (c *wrapper) String(text string) error {
	c.res.Header().Set("Content-Type", "text/plain")
	_, err := c.res.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

func (c *wrapper) Stream(contentType string, rd io.Reader) error {
	c.res.Header().Set("Content-Type", contentType)
	_, err := io.Copy(c.res, rd)
	return err
}

func (c *wrapper) Reset(res http.ResponseWriter, req *http.Request) {
	c.w.reset(res)
	c.res = res
	c.req = req
}

func (c *wrapper) Deadline() (time.Time, bool) {
	if c.req == nil {
		return time.Time{}, false
	}
	return c.req.Context().Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Done()
}

func (c *wrapper) Err() error {
	if c.req == nil {
		return context.Canceled
	}
	return c.req.Context().Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Value(key)
}

func (c *wrapper) Client() mqtt.Client {
	return ClientFromContext(c.req.Context())
}

type messageKey struct{}

func MessageFromContext(ctx context.Context) mqtt.Message {
	msg := ctx.Value(messageKey{})
	if msg == nil {
		return nil
	}
	return msg.(mqtt.Message)
}

func WithMessage(ctx context.Context, msg mqtt.Message) context.Context {
	if msg == nil {
		return ctx
	}
	return context.WithValue(ctx, messageKey{}, msg)
}

type clientKey struct{}

func ClientFromContext(ctx context.Context) mqtt.Client {
	msg := ctx.Value(clientKey{})
	if msg == nil {
		return nil
	}
	return msg.(mqtt.Client)
}

func WithClient(ctx context.Context, client mqtt.Client) context.Context {
	if client == nil {
		return ctx
	}
	return context.WithValue(ctx, clientKey{}, client)
}
