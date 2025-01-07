package mqtt

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(ctx context.Context, contentType string, in interface{}) (body []byte, err error)

// ClientOption is MQTT client option.
type ClientOption func(*clientOptions)

// Client is an MQTT transport client.
type clientOptions struct {
	ctx           context.Context
	encoder       EncodeRequestFunc
	publishMQTTFn PublishMQTTFn
	middleware    []middleware.Middleware
	timeout       time.Duration
}

// WithTimeout with client request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// WithRequestEncoder with client request encoder.
func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}
func WithPublishMQTTFn(fn PublishMQTTFn) ClientOption {
	return func(o *clientOptions) {
		o.publishMQTTFn = fn
	}
}

// Client is an MQTT client.
type Client struct {
	opts clientOptions
}

// NewClient returns an MQTT client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		ctx:     ctx,
		timeout: 2000 * time.Millisecond,
		encoder: DefaultRequestEncoder,
	}
	for _, o := range opts {
		o(&options)
	}
	return &Client{
		opts: options,
	}, nil
}

func (client *Client) Publish(ctx context.Context, topic string, qos byte, retain bool, args interface{}) error {
	var data []byte
	var err error
	if args != nil {
		data, err = client.opts.encoder(ctx, "json", args)
		if err != nil {
			return err
		}
	}
	h := func(ctx context.Context, _ interface{}) (interface{}, error) {
		if client.opts.publishMQTTFn == nil {
			log.Error("publishMQTTFn is nil")
			return nil, nil
		}
		err := client.opts.publishMQTTFn(topic, qos, retain, data)
		return nil, err
	}

	if len(client.opts.middleware) > 0 {
		h = middleware.Chain(client.opts.middleware...)(h)
	}
	_, err = h(ctx, args)
	return err
}

// Close tears down the MQTT client.
func (client *Client) Close() error {
	return nil
}

// DefaultRequestEncoder is an MQTT request encoder.
func DefaultRequestEncoder(_ context.Context, contentType string, in interface{}) ([]byte, error) {

	if contentType == "" {
		contentType = "json"
	}
	body, err := encoding.GetCodec(contentType).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, nil
}
