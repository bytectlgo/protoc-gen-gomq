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
	MQRouter.Handle("/test/:id", func(ctx context.Context, client any, msg any, ps crouter.Params) {
		fmt.Println(ps)
		fmt.Println(msg)
	})
	msg := &message{
		topic:   "/test/1234",
		payload: []byte("test"),
	}
	MQRouter.Serve(context.Background(), msg.topic, nil, msg)
}
