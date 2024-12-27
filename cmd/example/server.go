package main

import (
	"fmt"
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
	var subscribeMQTTFn mqtt.SubscribeMQTTFn
	opts := pmqtt.NewClientOptions()
	opts.OnConnectionLost = func(client pmqtt.Client, err error) {
		reader := client.OptionsReader()
		fmt.Printf("mqtt lost connect client id: %v\n", reader.ClientID())
	}
	opts.OnReconnecting = func(client pmqtt.Client, options *pmqtt.ClientOptions) {
		fmt.Printf("mqtt reconnecting client id: %v\n", options.ClientID)
	}
	opts.OnConnect = func(client pmqtt.Client) {
		fmt.Println("mqtt connected")
		// 定阅消息
		if subscribeMQTTFn == nil {
			fmt.Println("subscribeTopic is nil")
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
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	mid := mqtt.Middleware(
		recovery.Recovery(),
		validate.Validator(),
		logging.Server(log.NewStdLogger(os.Stdout)),
	)
	server := mqtt.NewServer(mqtt.WithClientOption(opts), mid)
	// 赋值定阅函数
	subscribeMQTTFn = server.MakeSubscribeMQTTFn()
	gencode.RegisterExampleMQServer(server, service)
	return server
}
