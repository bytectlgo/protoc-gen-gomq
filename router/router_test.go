package router

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytectlgo/crouter"
)

type message struct {
	topic   string
	payload []byte
}

func (m *message) Topic() string {
	return m.topic
}
func (m *message) Payload() []byte {
	return m.payload
}
func (m *message) MessageID() uint16 {
	return 0
}

func TestMQRouter(t *testing.T) {

	MQRouter := NewMQRouter()
	MQRouter.SaveMatchedRoutePath = false
	h := func(ctx context.Context, client any, msg any, ps crouter.Params) error {
		fmt.Println(ps)
		fmt.Println(msg)
		return nil
	}
	filter := func(next Handle) Handle {
		return func(ctx context.Context, client any, msg any, ps crouter.Params) error {
			fmt.Println("filter")
			return next(ctx, client, msg, ps)
		}
	}
	filter2 := func(next Handle) Handle {
		return func(ctx context.Context, client any, msg any, ps crouter.Params) error {
			fmt.Println("filter2")
			return next(ctx, client, msg, ps)
		}
	}
	MQRouter.Handle("/test/:id", h, filter, filter2)
	msg := &message{
		topic:   "/test/1234",
		payload: []byte("test"),
	}
	MQRouter.Serve(context.Background(), msg.topic, nil, msg)
}
