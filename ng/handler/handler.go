package handler

import (
	"reflect"
	"runtime"
	"strings"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

type Handler func(*context.Context) *response.Response

func Name(h Handler) string {
	s := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	s = s[strings.LastIndex(s, ".")+1:]
	s = strings.TrimSuffix(s, "Handler")
	return s
}
