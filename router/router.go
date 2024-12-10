package router

import (
	"context"

	"github.com/bytectlgo/crouter"
	"github.com/go-kratos/kratos/v2/log"
)

type Handle func(ctx context.Context, client any, msg any, ps crouter.Params) error
type FilterFunc func(next Handle) Handle

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(ctx context.Context, client any, msg any, ps crouter.Params) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(ctx context.Context, client any, msg any, ps crouter.Params, resp any) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(ctx context.Context, client any, msg any, ps crouter.Params, err error)

const (
	MQMethod = "MQ"
)

type option func(r *MQRouter)

func WithNotFound(h func(context.Context)) option {
	return func(r *MQRouter) {
		r.Router.NotFound = crouter.HandlerFunc(h)
	}
}

func WithErrorEncoder(ene EncodeErrorFunc) option {
	return func(r *MQRouter) {
		r.ene = ene
	}
}
func WithFilters(filters ...FilterFunc) option {
	return func(r *MQRouter) {
		r.filters = filters
	}
}

type MQRouter struct {
	*crouter.Router
	client  any
	filters []FilterFunc
	ene     EncodeErrorFunc
}

func NewMQRouter(opts ...option) *MQRouter {
	mqr := &MQRouter{
		Router:  crouter.New(),
		filters: []FilterFunc{},
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
	for _, opt := range opts {
		opt(mqr)
	}
	return mqr
}

func (r *MQRouter) SetErrorEncoder(ene EncodeErrorFunc) {
	r.ene = ene
}

func (r *MQRouter) Use(filters ...FilterFunc) {
	r.filters = append(r.filters, filters...)
}

func (r *MQRouter) Handle(path string, h Handle, filters ...FilterFunc) {
	next := FilterChain(filters...)(h)
	next = FilterChain(r.filters...)(next)
	tnext := func(ctx context.Context, ps crouter.Params) {
		rctx, ok := ctx.(*wrapper)
		if !ok {
			log.Warnf("ctx is not a wrapper")
			return
		}
		msg := rctx.msg
		client := rctx.client
		if err := next(ctx, client, msg, ps); err != nil {
			if r.ene != nil {
				r.ene(ctx, client, msg, ps, err)
			} else {
				log.Error("Handle error: ", err)
			}
		}
	}
	if _, _, ok := r.Router.Lookup(MQMethod, path); ok {
		return
	}
	r.Router.Handle(MQMethod, path, tnext)
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
