package handler

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

var contextType = reflect.TypeOf(&context.Context{})

func MustFunc(function interface{}, pars ...*param.Param) Handler {
	h, err := Func(function, pars...)
	if err != nil {
		panic(err)
	}
	return h
}

func Func(function interface{}, pars ...*param.Param) (Handler, error) {
	fun := reflect.ValueOf(function)
	funType := fun.Type()

	if fun.Kind() != reflect.Func {
		return nil, errors.Errorf("not a function: %v", function)
	}

	funName := runtime.FuncForPC(fun.Pointer()).Name()

	if funType.NumOut() != 1 {
		return nil, errors.Errorf("function %s must return 1 value", funName)
	}

	cardinality := len(pars)

	if cardinality+1 != funType.NumIn() {
		s := ""
		if cardinality+1 != 1 {
			s = "s"
		}
		return nil, errors.Errorf("function %s must have %d parameter%s of type %s", funName, cardinality+1, s, types(pars))
	}

	if funType.In(0) != contextType {
		return nil, errors.Errorf("function %s parameter 0 must be a %s", funName, contextType)
	}

	for i, par := range pars {
		t := funType.In(i + 1)
		if !t.AssignableTo(par.Type) {
			return nil, errors.Errorf("function %s parameter %d must be %s (assignable to %s type) but is %s", funName, i+1, par.Name, par.Type, t)
		}
	}

	return func(ctx *context.Context) *response.Response {
		args := ctx.Stub.GetArgs()
		argc := len(args) - 1
		if argc > 0 {
			args = args[1:]
		} else {
			args = nil
		}
		if err := ValidateArgCount(ctx, cardinality); err != nil {
			return response.BadRequest(err.Error())
		}
		vals := []reflect.Value{reflect.ValueOf(ctx)}
		argi := 0
		for i, par := range pars {
			v, c, err := par.Parse(args[argi:])
			if err != nil {
				return response.BadRequest("%s argument %d: %v", par.Name, i+1, err)
			}
			argi += c
			vals = append(vals, reflect.ValueOf(v))
		}
		ret := fun.Call(vals)[0].Interface()
		if res, ok := ret.(*response.Response); ok {
			return res
		}
		return response.OK(ret)
	}, nil

}

func types(pars []*param.Param) string {
	ss := []string{contextType.String()}
	for _, par := range pars {
		ss = append(ss, par.Type.String())
	}
	return strings.Join(ss, ", ")
}