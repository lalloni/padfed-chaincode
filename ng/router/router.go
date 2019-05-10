package router

import (
	"fmt"
	"strings"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
)

type Name string

func (n *Name) String() string {
	return strings.ToLower(string(*n)) // para que los nombres de las funciones sean case-insensitive
}

type Router interface {
	InitHandler() handler.Handler
	SetInitHandler(authorization.Check, handler.Handler)
	Handler(Name) handler.Handler
	SetHandler(Name, authorization.Check, handler.Handler)
}

func New(c *Config) Router {
	r := &router{
		functionHandlers: map[string]handler.Handler{},
	}
	if c != nil {
		if c.Init != nil {
			if c.Init.Check != nil || c.Init.Handler != nil {
				r.SetInitHandler(c.Init.Check, HandlerDefault(c.Init.Handler, handler.SuccessHandler))
			}
		}
		for _, fun := range c.Funs {
			r.SetHandler(NameDefault(fun.Name, fun.Handler), fun.Check, fun.Handler)
		}
	}
	return r
}

type router struct {
	initHandler      handler.Handler
	functionHandlers map[string]handler.Handler
}

func (r *router) InitHandler() handler.Handler {
	return r.initHandler
}

func (r *router) SetInitHandler(ch authorization.Check, h handler.Handler) {
	if ch != nil {
		r.initHandler = authorization.Handler("init", ch, h)
	} else {
		r.initHandler = h
	}
}

func (r *router) Handler(n Name) handler.Handler {
	return r.functionHandlers[n.String()]
}

func (r *router) SetHandler(n Name, ch authorization.Check, h handler.Handler) {
	if ch != nil {
		r.functionHandlers[n.String()] = authorization.Handler(fmt.Sprintf("invoke function %q", n), ch, h)
	} else {
		r.functionHandlers[n.String()] = h
	}
}
