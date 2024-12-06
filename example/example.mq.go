package example

import (
	"github.com/bytectlgo/protoc-gen-gomq/mq"
	"github.com/go-kratos/kratos/v2/log"
)

type ExampleMQServer interface {
	DeviceInfo(mq.Context, *DeviceInfoPulish) (*DeviceInfoReply, error)
	DeviceProperty(mq.Context, *DevicePropertyulish) (*DevicePropertyReply, error)
}

func SubscribeExampleMQServer(c mq.Client, m *mq.MQSubscribe) {
	m.Subscribe(c, "$share/device//$SYS/brokers/:node/clients/:clientId/connected", 0)
	m.Subscribe(c, "$share/device//device/:device_id/property", 0)
}
func RegisterExampleMQServer(s *mq.Server, srv ExampleMQServer) {
	r := s.Route()
	r.Handle("$share/device//$SYS/brokers/:node/clients/:clientId/connected", _DeviceInfo_DeviceInfoMQ_Handler(srv))
	r.Handle("$share/device//device/:device_id/property", _DeviceProperty_DevicePropertyMQ_Handler(srv))
}
func _ExampleMQServer_DeviceInfoMQ_Handler(srv ExampleMQServer) func(mq.Context) {
	return func(ctx mq.Context) {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &DeviceInfoPulish{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return
		}
		log.Debugf("receive mq topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return
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
			ctx.ReplyErr(err)
			return
		}
		// reply topic:/$SYS/brokers/:node/clients/:clientId/connected_reply
		err = ctx.Reply(reply)
		if err != nil {
			log.Error("DeviceInfo error:", err)
			ctx.ReplyErr(err)
			return
		}
	}
}
func _ExampleMQServer_DevicePropertyMQ_Handler(srv ExampleMQServer) func(mq.Context) {
	return func(ctx mq.Context) {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &DevicePropertyulish{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return
		}
		log.Debugf("receive mq topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return
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
			ctx.ReplyErr(err)
			return
		}
		// reply topic:/device/:device_id/property_reply
		err = ctx.Reply(reply)
		if err != nil {
			log.Error("DeviceProperty error:", err)
			ctx.ReplyErr(err)
			return
		}
	}
}
