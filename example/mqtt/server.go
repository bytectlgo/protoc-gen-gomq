package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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

func main() {
	opts := pmqtt.NewClientOptions()
	opts.OnConnectionLost = func(client pmqtt.Client, err error) {
		fmt.Println("Connection lost to MQTT server")
	}
	opts.OnReconnecting = func(client pmqtt.Client, options *pmqtt.ClientOptions) {
		fmt.Println("Reconnecting to MQTT server")
	}
	opts.OnConnect = func(client pmqtt.Client) {
		fmt.Println("Connected to MQTT server")
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
	server.Start(context.Background())
	topic := "/test/{product}/{device}"
	subscribeTopic := mqtt.MakeSubscribeFn(server)
	subscribeTopic(topic, 1)
	route := server.Route("/")
	route.POST(topic, info)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	server.Stop(context.Background())
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
