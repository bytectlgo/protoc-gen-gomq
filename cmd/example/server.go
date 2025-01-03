package main

import (
	"os"

	"example/gencode"

	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
)

func NewMQTTSever(
	service *Service,
) *mqtt.Server {
	glogger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.Caller(4),
	)
	log.SetLogger(glogger)
	var subscribeMQTTFn mqtt.SubscribeMQTTFn
	opts := pmqtt.NewClientOptions()
	opts.OnConnectionLost = func(client pmqtt.Client, err error) {
		reader := client.OptionsReader()
		log.Debugf("mqtt lost connect client id: %v", reader.ClientID())
	}
	opts.OnReconnecting = func(client pmqtt.Client, options *pmqtt.ClientOptions) {
		log.Debugf("mqtt reconnecting client id: %v", options.ClientID)
	}
	opts.OnConnect = func(client pmqtt.Client) {
		log.Debugf("mqtt connected")
		// 定阅消息
		if subscribeMQTTFn == nil {
			log.Debugf("subscribeTopic is nil")
			return
		}
		gencode.SubscribeExampleMQServer(subscribeMQTTFn)
	}
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("server")
	opts.SetUsername("user")
	opts.SetPassword("password")
	opts.SetResumeSubs(true)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(func(client pmqtt.Client, msg pmqtt.Message) {
		log.Debugf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	})
	mid := mqtt.Middleware(
		recovery.Recovery(),
		logging.Server(glogger),
		validate.Validator(),
	)
	server := mqtt.NewServer(mqtt.WithClientOption(opts), mid)
	// 赋值定阅函数
	subscribeMQTTFn = server.MakeSubscribeMQTTFn()
	gencode.RegisterExampleMQServer(server, service)
	return server
}
