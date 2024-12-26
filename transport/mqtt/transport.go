package mqtt

import (
	"context"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ Transporter = (*Transport)(nil)

// Transporter is mqtt Transporter
type Transporter interface {
	transport.Transporter
	Request() *http.Request
	Client() *mqtt.Client
}

// Transport is an MQTT transport.
type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
	request     *http.Request
	response    http.ResponseWriter
	client      *mqtt.Client
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return "mqtt"
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Request returns the MQTT request.
func (tr *Transport) Request() *http.Request {
	return tr.request
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

// Client returns the MQTT client.
func (tr *Transport) Client() *mqtt.Client {
	return tr.client
}

// SetOperation sets the transport operation.
func SetOperation(ctx context.Context, op string) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			tr.operation = op
		}
	}
}

// RequestFromServerContext returns request from context.
func RequestFromServerContext(ctx context.Context) (*http.Request, bool) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			return tr.request, true
		}
	}
	return nil, false
}

// ClientFromServerContext returns client from context.
func ClientFromServerContext(ctx context.Context) (*mqtt.Client, bool) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			return tr.client, true
		}
	}
	return nil, false
}

type headerCarrier http.Header

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Add append value to key-values pair.
func (hc headerCarrier) Add(key string, value string) {
	http.Header(hc).Add(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}
