package mqtt

import (
	"context"
	"strings"

	"github.com/bytectlgo/crouter"
	"github.com/bytectlgo/protoc-gen-gomq/router"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

type Handle func(ctx context.Context, client pmqtt.Client, msg pmqtt.Message, ps crouter.Params)

type MQTTMsgServer struct {
	mqttClient pmqtt.Client
	router     *router.MQRouter
}

func NewMQTTMsgServer(client pmqtt.Client) *MQTTMsgServer {
	srv := &MQTTMsgServer{
		mqttClient: client,
	}
	srv.router = router.NewMQRouter()
	return srv
}

func (s *MQTTMsgServer) Subscribe(topic string, qos byte, h Handle) error {
	subscribeTopic := getSubscribeTopic(topic)
	routePath := getRoutePathFromTopic(topic)
	s.mqttClient.Subscribe(subscribeTopic, qos, s.serve)
	hnext := func(ctx context.Context, client any, msg any, ps crouter.Params) {
		c, ok := client.(pmqtt.Client)
		if !ok {
			log.Errorf("client is not a pmqtt.Client")
			return
		}
		m, ok := msg.(pmqtt.Message)
		if !ok {
			log.Errorf("msg is not a pmqtt.Message")
			return
		}
		h(ctx, c, m, ps)
	}
	s.router.Handle(routePath, hnext)
	return nil
}
func (s *MQTTMsgServer) serve(client pmqtt.Client, msg pmqtt.Message) {
	path := getRoutePathFromTopic(msg.Topic())
	s.router.Serve(context.Background(), path, client, msg)
}

func getRoutePathFromTopic(topic string) string {
	if strings.HasPrefix(topic, "/") {
		return topic
	}
	return "/noSlashRoot/" + topic
}

func getSubscribeTopic(topic string) string {
	dirs := strings.Split(topic, "/")
	for i, dir := range dirs {
		if dir == "" {
			continue
		}
		if dir[0] == ':' {
			dirs[i] = "+"
		}
		if dir[0] == '*' {
			dirs[i] = "#"
		}
	}
	return strings.Join(dirs, "/")
}
