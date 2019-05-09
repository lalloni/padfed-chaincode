package router

import (
	"fmt"
	"reflect"
	"runtime"
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
	SetHandlerFunc(authorization.Check, handler.Handler)
}

func New() Router {
	return &router{
		functionHandlers: map[string]handler.Handler{},
	}
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
		r.initHandler = handler.AuthorizationHandler("init", ch, h)
	} else {
		r.initHandler = h
	}
}

func (r *router) Handler(n Name) handler.Handler {
	return r.functionHandlers[n.String()]
}

func (r *router) SetHandler(n Name, ch authorization.Check, h handler.Handler) {
	if ch != nil {
		r.functionHandlers[n.String()] = handler.AuthorizationHandler(fmt.Sprintf("invoke function %q", n), ch, h)
	} else {
		r.functionHandlers[n.String()] = h
	}
}

func (r *router) SetHandlerFunc(ch authorization.Check, h handler.Handler) {
	s := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	s = s[strings.LastIndex(s, ".")+1:]
	s = strings.TrimSuffix(s, "Handler")
	r.SetHandler(Name(s), ch, h)
}
