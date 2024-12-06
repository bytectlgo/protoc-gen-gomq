package mq

import "context"

type Context interface {
	context.Context
	Bind(v interface{}) error
	BindVars(v interface{}) error
}
