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
	subscribeMQTTFn("$share/device//device/{deviceKey}/service/{action}", 1)
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.EventPost(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("EventPost:", err)
			}
		}
		if err != nil {
			log.Error("EventPost:", err)
			return err
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "0")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ServiceRequest(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("ServiceRequest:", err)
			}
		}
		if err != nil {
			log.Error("ServiceRequest:", err)
			return err
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "0")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ServiceReply(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("ServiceReply:", err)
			}
		}
		if err != nil {
			log.Error("ServiceReply:", err)
			return err
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "true")
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
func ClientSubscribeExampleMQServer(subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	subscribeMQTTFn("/device/{deviceKey}/event/{action}/post_reply", 0)
	subscribeMQTTFn("/device/{deviceKey}/service/{action}_reply", 0)
	subscribeMQTTFn("/device/{deviceKey}/service/{action}_reply", 1)
}
func ClientRegisterExampleMQServer(s *mqtt.Server, srv ClientExampleMQServer) {
	r := s.Route("/")
	r.POST("/device/{deviceKey}/event/{action}/post_reply", _ClientExampleMQServer_EventPostMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}_reply", _ClientExampleMQServer_ServiceRequestMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}_reply", _ClientExampleMQServer_ServiceReplyMQ_Handler(srv))
}

type ClientExampleMQServer interface {
	ClientEventPost(context.Context, *Reply) error
	ClientServiceRequest(context.Context, *Reply) error
	ClientServiceReply(context.Context, *Reply) error
}

func _ClientExampleMQServer_EventPostMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &Reply{}
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientEventPost(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		return err
	}
}
func _ClientExampleMQServer_ServiceRequestMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &Reply{}
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientServiceRequest(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		return err
	}
}
func _ClientExampleMQServer_ServiceReplyMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &Reply{}
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
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientServiceReply(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		return err
	}
}

type ClientExampleMQServerImpl struct {
	client *mqtt.Client
}

func NewClientExampleMQServerImpl(client *mqtt.Client) *ClientExampleMQServerImpl {
	return &ClientExampleMQServerImpl{
		client: client,
	}
}
func (c *ClientExampleMQServerImpl) EventPost(ctx context.Context, in *ThingReq) error {
	topic := "/device/{deviceKey}/event/{action}/post"
	path := binding.EncodeURL(topic, in, false)
	return c.client.Publish(ctx, path, 0, false, in)
}
func (c *ClientExampleMQServerImpl) ServiceRequest(ctx context.Context, in *ThingReq) error {
	topic := "/device/{deviceKey}/service/{action}"
	path := binding.EncodeURL(topic, in, false)
	return c.client.Publish(ctx, path, 0, false, in)
}
func (c *ClientExampleMQServerImpl) ServiceReply(ctx context.Context, in *ThingReq) error {
	topic := "/device/{deviceKey}/service/{action}"
	path := binding.EncodeURL(topic, in, false)
	return c.client.Publish(ctx, path, 1, true, in)
}
