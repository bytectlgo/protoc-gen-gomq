package gencode

import (
	"context"

	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

type ExampleMQServer interface {
	EventPost(context.Context, *ThingReq) (*Reply, error)
	ServiceRequest(context.Context, *ThingReq) (*Reply, error)
	ServiceReply(context.Context, *ThingReq) (*Reply, error)
}

func SubscribeExampleMQServer(subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	subscribeMQTTFn("$share/device//device/{deviceKey}/event/{action}/post", 0)
	subscribeMQTTFn("$share/device//device/{deviceKey}/service/{action}", 0)
	subscribeMQTTFn("$share/device//device/{deviceKey}/service/{action}", 0)
}
func RegisterExampleMQServer(s *mqtt.Server, srv ExampleMQServer) {
	r := s.Route("/")
	r.POST("/device/{deviceKey}/event/{action}/post", _ExampleMQServer_EventPostMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}", _ExampleMQServer_ServiceRequestMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}", _ExampleMQServer_ServiceReplyMQ_Handler(srv))
}
func _ExampleMQServer_EventPostMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		err = ctx.BindVars(in)
		if err != nil {
			log.Error("bind vars error:", err)
			return err
		}
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
		}
		log.Debugf("receive mq request:%+v", in)
		reply, err := srv.EventPost(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("EventPost error:", err)
			}
		}
		if err != nil {
			log.Error("EventPost error:", err)
			return err
		}
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
		pattern := "/device/{deviceKey}/event/{action}/post_reply"
		topic := binding.EncodeURL(pattern, in, false)
		err = ctx.JSON(topic, reply)
		if err != nil {
			log.Error("EventPost error:", err)
			return err
		}
		return nil
	}
}
func _ExampleMQServer_ServiceRequestMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		err = ctx.BindVars(in)
		if err != nil {
			log.Error("bind vars error:", err)
			return err
		}
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
		}
		log.Debugf("receive mq request:%+v", in)
		reply, err := srv.ServiceRequest(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("ServiceRequest error:", err)
			}
		}
		if err != nil {
			log.Error("ServiceRequest error:", err)
			return err
		}
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
		pattern := "/device/{deviceKey}/service/{action}_reply"
		topic := binding.EncodeURL(pattern, in, false)
		err = ctx.JSON(topic, reply)
		if err != nil {
			log.Error("ServiceRequest error:", err)
			return err
		}
		return nil
	}
}
func _ExampleMQServer_ServiceReplyMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		err = ctx.BindVars(in)
		if err != nil {
			log.Error("bind vars error:", err)
			return err
		}
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
		}
		log.Debugf("receive mq request:%+v", in)
		reply, err := srv.ServiceReply(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("ServiceReply error:", err)
			}
		}
		if err != nil {
			log.Error("ServiceReply error:", err)
			return err
		}
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
		// ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
		pattern := "/device/{deviceKey}/service/{action}_reply"
		topic := binding.EncodeURL(pattern, in, false)
		err = ctx.JSON(topic, reply)
		if err != nil {
			log.Error("ServiceReply error:", err)
			return err
		}
		return nil
	}
}
