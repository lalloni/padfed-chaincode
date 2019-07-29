package organizaciones

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
)

func AddHandlers(r router.Router) {
	r.SetHandler("GetOrganizacion", common.Free, getHandler)
	r.SetHandler("GetOrganizacionAll", common.Free, getAllHandler)
}

func getHandler(ctx *context.Context) *response.Response {
	var id uint64
	_, err := handler.ExtractArgs(ctx.Stub.GetArgs()[1:], param.Uint64Var(&id))
	if err != nil {
		return response.BadRequest(err.Error())
	}
	org := GetByID(id)
	if org == nil {
		return response.NotFoundWithMessage("organizacion identified with %v not found", id)
	}
	return response.OK(org)
}

func getAllHandler(ctx *context.Context) *response.Response {
	_, err := handler.ExtractArgs(ctx.Stub.GetArgs()[1:]) // no parameters
	if err != nil {
		return response.BadRequest(err.Error())
	}
	return response.OK(GetAll())
}
