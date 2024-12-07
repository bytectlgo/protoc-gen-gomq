package router

import (
	"context"
	"time"

	"github.com/bytectlgo/crouter"
)

type Context interface {
	context.Context
	crouter.RouterInfo
}

var _ Context = (*wrapper)(nil)

type wrapper struct {
	ctx    context.Context
	client any
	msg    any
	ps     *crouter.Params
	path   string
	method string
}

func (c *wrapper) GetPath() string {
	return c.path
}

func (c *wrapper) Method() string {
	return c.method
}

func (c *wrapper) Deadline() (time.Time, bool) {
	if c.ctx == nil {
		return time.Time{}, false
	}
	return c.ctx.Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Done()
}

func (c *wrapper) Err() error {
	if c.ctx == nil {
		return context.Canceled
	}
	return c.ctx.Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Value(key)
}
