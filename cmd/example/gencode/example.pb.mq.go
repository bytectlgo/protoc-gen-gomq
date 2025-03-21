package gencode

import (
	"context"
	"fmt"

	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

const OperationEventPost = "/gencode.ExampleMQServer/EventPost"
const OperationServiceRequest = "/gencode.ExampleMQServer/ServiceRequest"
const OperationServiceReply = "/gencode.ExampleMQServer/ServiceReply"

type ExampleMQServer interface {
	// 设备属性,事件上报
	EventPost(context.Context, *ThingReq) (*Reply, error)
	// 服务指令下发
	// 服务器下发指令给设备
	// 设备回复指令执行结果
	ServiceRequest(context.Context, *ThingReq) (*Reply, error)
	// 服务指令回复
	ServiceReply(context.Context, *ThingReq) (*Reply, error)
}

func SubscribeExampleMQServer(groupPrefix string, subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/event/{action}/post", 0)
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/service/{action}", 0)
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/service/{action}", 1)
}
func RegisterExampleMQServer(s *mqtt.Server, srv ExampleMQServer) {
	r := s.Route("/")
	r.POST("/device/{deviceKey}/event/{action}/post", _ExampleMQServer_EventPostMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}", _ExampleMQServer_ServiceRequestMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}", _ExampleMQServer_ServiceReplyMQ_Handler(srv))
}
func _ExampleMQServer_EventPostMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}

		mqtt.SetOperation(ctx, OperationEventPost)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.EventPost(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "0")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
		pattern := "/device/{deviceKey}/event/{action}/post_reply"
		topic := binding.EncodeURL(pattern, in, false)
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_TOPIC_HEADER, topic)
		err = ctx.JSON(reply)
		if err != nil {
			return fmt.Errorf("json error:%v", err)
		}
		return nil
	}
}
func _ExampleMQServer_ServiceRequestMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}

		mqtt.SetOperation(ctx, OperationServiceRequest)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ServiceRequest(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "0")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
		pattern := "/device/{deviceKey}/service/{action}_reply"
		topic := binding.EncodeURL(pattern, in, false)
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_TOPIC_HEADER, topic)
		err = ctx.JSON(reply)
		if err != nil {
			return fmt.Errorf("json error:%v", err)
		}
		return nil
	}
}
func _ExampleMQServer_ServiceReplyMQ_Handler(srv ExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &ThingReq{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}

		mqtt.SetOperation(ctx, OperationServiceReply)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ServiceReply(ctx, req.(*ThingReq))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "true")
		pattern := "/device/{deviceKey}/service/{action}_reply"
		topic := binding.EncodeURL(pattern, in, false)
		ctx.Response().Header().Set(mqtt.MQTT_REPLY_TOPIC_HEADER, topic)
		err = ctx.JSON(reply)
		if err != nil {
			return fmt.Errorf("json error:%v", err)
		}
		return nil
	}
}
func ClientSubscribeExampleMQServer(groupPrefix string, subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/event/{action}/post_reply", 0)
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/service/{action}_reply", 0)
	subscribeMQTTFn(groupPrefix+"/device/{deviceKey}/service/{action}_reply", 1)
}
func ClientRegisterExampleMQServer(s *mqtt.Server, srv ClientExampleMQServer) {
	r := s.Route("/")
	r.POST("/device/{deviceKey}/event/{action}/post_reply", _ClientExampleMQServer_EventPostMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}_reply", _ClientExampleMQServer_ServiceRequestMQ_Handler(srv))
	r.POST("/device/{deviceKey}/service/{action}_reply", _ClientExampleMQServer_ServiceReplyMQ_Handler(srv))
}

type ClientExampleMQServer interface {
	// 设备属性,事件上报
	ClientEventPost(context.Context, *Reply) error
	// 服务指令下发
	// 服务器下发指令给设备
	// 设备回复指令执行结果
	ClientServiceRequest(context.Context, *Reply) error
	// 服务指令回复
	ClientServiceReply(context.Context, *Reply) error
}

func _ClientExampleMQServer_EventPostMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &Reply{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}
		mqtt.SetOperation(ctx, OperationEventPost)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientEventPost(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		return nil
	}
}
func _ClientExampleMQServer_ServiceRequestMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &Reply{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}
		mqtt.SetOperation(ctx, OperationServiceRequest)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientServiceRequest(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		return nil
	}
}
func _ClientExampleMQServer_ServiceReplyMQ_Handler(srv ClientExampleMQServer) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &Reply{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}
		mqtt.SetOperation(ctx, OperationServiceReply)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.ClientServiceReply(ctx, req.(*Reply))
			return nil, err
		})
		_, err = h(ctx, in)
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		return nil
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

// 设备属性,事件上报
func (c *ClientExampleMQServerImpl) EventPost(ctx context.Context, in *ThingReq, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/event/{action}/post"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationEventPost))
	return c.client.Publish(ctx, path, 0, false, in, opts...)
}

// 服务指令下发
// 服务器下发指令给设备
// 设备回复指令执行结果
func (c *ClientExampleMQServerImpl) ServiceRequest(ctx context.Context, in *ThingReq, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/service/{action}"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationServiceRequest))
	return c.client.Publish(ctx, path, 0, false, in, opts...)
}

// 服务指令回复
func (c *ClientExampleMQServerImpl) ServiceReply(ctx context.Context, in *ThingReq, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/service/{action}"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationServiceReply))
	return c.client.Publish(ctx, path, 1, true, in, opts...)
}

type ClientReplyExampleMQServerImpl struct {
	client *mqtt.Client
}

func NewClientReplyExampleMQServerImpl(client *mqtt.Client) *ClientReplyExampleMQServerImpl {
	return &ClientReplyExampleMQServerImpl{
		client: client,
	}
}

// 设备属性,事件上报
func (c *ClientReplyExampleMQServerImpl) EventPost(ctx context.Context, in *Reply, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/event/{action}/post_reply"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationEventPost))
	return c.client.Publish(ctx, path, 0, false, in, opts...)
}

// 服务指令下发
// 服务器下发指令给设备
// 设备回复指令执行结果
func (c *ClientReplyExampleMQServerImpl) ServiceRequest(ctx context.Context, in *Reply, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/service/{action}_reply"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationServiceRequest))
	return c.client.Publish(ctx, path, 0, false, in, opts...)
}

// 服务指令回复
func (c *ClientReplyExampleMQServerImpl) ServiceReply(ctx context.Context, in *Reply, opts ...mqtt.CallOption) error {
	topic := "/device/{deviceKey}/service/{action}_reply"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(OperationServiceReply))
	return c.client.Publish(ctx, path, 1, true, in, opts...)
}
