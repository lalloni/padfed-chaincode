package router

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func FunctionsHandler(r Router) handler.Handler {
	return func(ctx *context.Context) *response.Response {
		if err := handler.CheckArgsCount(ctx, 0); err != nil {
			return response.BadRequest(err.Error())
		}
		return response.OK(r.Functions())
	}
}
