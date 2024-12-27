package example

import (
	"context"

	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
	"github.com/go-kratos/kratos/v2/log"
)

type ExampleMQServer interface {
	DeviceInfo(context.Context, *DeviceInfoPulish) (*DeviceInfoReply, error)
	DeviceProperty(context.Context, *DevicePropertyulish) (*DevicePropertyReply, error)
}

func SubscribeExampleMQServer(subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	subscribeMQTTFn("$share/device//$SYS/brokers/:node/clients/:clientId/connected", 0)
	subscribeMQTTFn("$share/device//device/:device_id/property", 0)
}
func RegisterExampleMQServer(s *mqtt.Server, srv ExampleMQServer) {
	r := s.Route("/")
	r.POST("$share/device//$SYS/brokers/:node/clients/:clientId/connected", _ExampleMQServer_DeviceInfoMQ_Handler(srv))
	r.POST("$share/device//device/:device_id/property", _ExampleMQServer_DevicePropertyMQ_Handler(srv))
}
func _ExampleMQServer_DeviceInfoMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &DeviceInfoPulish{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		log.Debugf("receive mq topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
		}
		log.Debugf("receive mq request:%+v", in)
		reply, err := srv.DeviceInfo(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("DeviceInfo error:", err)
			}
		}
		if err != nil {
			log.Error("DeviceInfo error:", err)
			return err
		}
		// reply topic:/$SYS/brokers/:node/clients/:clientId/connected_reply
		err = ctx.Reply(reply)
		if err != nil {
			log.Error("DeviceInfo error:", err)
			return err
		}
	}
}
func _ExampleMQServer_DevicePropertyMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &DevicePropertyulish{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		log.Debugf("receive mq topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
		}
		log.Debugf("receive mq request:%+v", in)
		reply, err := srv.DeviceProperty(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("DeviceProperty error:", err)
			}
		}
		if err != nil {
			log.Error("DeviceProperty error:", err)
			return err
		}
		// reply topic:/device/:device_id/property_reply
		err = ctx.Reply(reply)
		if err != nil {
			log.Error("DeviceProperty error:", err)
			return err
		}
	}
}
