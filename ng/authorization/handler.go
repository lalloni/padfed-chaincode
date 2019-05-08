package authorization

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func Handler(action string, check Check, handler handler.Handler) handler.Handler {
	return func(ctx *context.Context) *response.Response {
		err := check(ctx)
		if err != nil {
			return response.Forbidden("%s forbidden: %s", action, err)
		}
		return handler(ctx)
	}
}
