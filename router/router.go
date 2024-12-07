package router

import (
	"context"

	"github.com/bytectlgo/crouter"
	"github.com/go-kratos/kratos/v2/log"
)

type Handle func(ctx context.Context, client any, msg any, ps crouter.Params)

const (
	MQMethod = "MQ"
)

type MQRouter struct {
	*crouter.Router
	client any
}

func NewMQRouter() *MQRouter {
	mqr := &MQRouter{
		Router: crouter.New(),
	}
	mqr.Router.NotFound = crouter.HandlerFunc(
		func(ctx context.Context) {
			rctx, ok := ctx.(Context)
			if !ok {
				log.Warnf("ctx is not a Context")
				return
			}
			log.Error("not found handler path: ", rctx.GetPath())
		},
	)
	return mqr
}

func (r *MQRouter) Handle(path string, h Handle) {
	th := func(ctx context.Context, ps crouter.Params) {
		rctx, ok := ctx.(*wrapper)
		if !ok {
			log.Warnf("ctx is not a wrapper")
			return
		}
		msg := rctx.msg
		client := rctx.client
		h(ctx, client, msg, ps)
	}
	if _, _, ok := r.Router.Lookup(MQMethod, path); ok {
		return
	}
	r.Router.Handle(MQMethod, path, th)
}

func (r *MQRouter) Serve(ctx context.Context, path string, client any, msg any) {
	rctx := &wrapper{
		ctx:    ctx,
		client: client,
		msg:    msg,
		path:   path,
		method: MQMethod,
	}
	r.Router.Serve(rctx)
}
