package router

import (
	"fmt"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
)

type Router interface {
	InitHandler() handler.Handler
	InvokeHandler(context.Function) handler.Handler
	SetInitHandler(authorization.Check, handler.Handler)
	SetInvokeHandler(authorization.Check, context.Function, handler.Handler)
}

func New() Router {
	return &router{}
}

type router struct {
	initHandler      handler.Handler
	functionHandlers map[context.Function]handler.Handler
}

func (r *router) InitHandler() handler.Handler {
	return r.initHandler
}

func (r *router) InvokeHandler(fn context.Function) handler.Handler {
	return r.functionHandlers[fn]
}

func (r *router) SetInitHandler(ch authorization.Check, h handler.Handler) {
	if ch != nil {
		r.initHandler = handler.AuthorizationHandler("init", ch, h)
	} else {
		r.initHandler = h
	}
}

func (r *router) SetInvokeHandler(ch authorization.Check, fn context.Function, h handler.Handler) {
	if ch != nil {
		r.functionHandlers[fn] = handler.AuthorizationHandler(fmt.Sprintf("invoke function %q", fn), ch, h)
	} else {
		r.functionHandlers[fn] = h
	}
}
