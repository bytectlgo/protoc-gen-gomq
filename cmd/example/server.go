package main

import (
	"fmt"

	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

type ReqInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
	Product  string `json:"product"`
	Device   string `json:"device"`
}

func NewMQTTSever(
// exampleMQServer ExampleMQServer,
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
		// SubscribeExampleMQServer(subscribeMQTTFn)
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
	server := mqtt.NewServer(mqtt.WithClientOption(opts))
	// 赋值定阅函数
	subscribeMQTTFn = server.MakeSubscribeMQTTFn()
	/// 注册路由
	//route := server.Route("/")
	// route.POST(topic, info)

	return server
}

func info(ctx mqtt.Context) error {
	msg := ctx.Message()
	log.Infof("message: %+v", msg)
	req := &ReqInfo{}
	err := ctx.Bind(req)
	if err != nil {
		log.Errorf("failed to bind request: %v", err)
		return err
	}
	err = ctx.BindVars(req)
	if err != nil {
		log.Errorf("failed to bind vars: %v", err)
		return err
	}

	log.Infof("req: %+v", req)
	replyTopic := msg.Topic() + "/reply"
	ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
	ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
	ctx.JSON(replyTopic, &ReqInfo{
		Username: "use2222",
		Password: "password2222",
		Version:  "1.0.0",
		Product:  "product2222",
		Device:   "device2222",
	})
	return nil
}
