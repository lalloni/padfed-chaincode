package handler

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func AuthorizationHandler(action string, check authorization.Check, handler Handler) Handler {
	return func(ctx *context.Context) *response.Response {
		err := check(ctx)
		if err != nil {
			return response.Forbidden("%s forbidden: %s", action, err)
		}
		return handler(ctx)
	}
}
