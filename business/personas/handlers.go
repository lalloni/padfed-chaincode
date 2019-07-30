package personas

import (
	"strconv"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handlerutil/crud"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/state"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestingHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(
		crud.Defaults, // i.e. get, getrange, has, put, putlist, del, delrange
		crud.WithIDParam(CUITParam),
		crud.WithItemParam(PersonaParam),
		crud.WithListParam(PersonaListParam),
		crud.WithValidator(validatePersona),
	)
	if !testing {
		opts = append(opts, crud.WithWriteCheck(common.AFIP))
	}
	crud.AddHandlers(r, Schema, opts...)

	r.SetHandler("QueryPersona", common.Free, QueryPersonaHandler)
	r.SetHandler("QueryPersonaBasica", common.Free, QueryPersonaBasicaHandler)
}

func validatePersona(ctx *context.Context, v interface{}) *response.Response {

	per := v.(*Persona)

	if per.ID == 0 {
		return response.BadRequest("id required")
	}

	if per.Persona == nil {
		exist, err := ctx.Store.HasComposite(Schema, per.ID)
		if err != nil {
			return response.Error("checking persona existence: %v", err)
		}
		if !exist {
			return response.BadRequest("persona is required when putting a new instance")
		}
	} else if per.Persona.ID != per.ID {
		return response.BadRequest("id %q and id %q must be equal", per.ID, per.Persona.ID)
	}

	return nil
}

func QueryPersonaHandler(ctx *context.Context) *response.Response {
	args, err := handler.ExtractArgs(ctx.Stub.GetArgs()[1:], CUITParam)
	if err != nil {
		return response.BadRequest("invalid argument: %v", err)
	}

	prefix := "per:" + strconv.FormatUint(args[0].(uint64), 10) + "#"
	query := state.Single(state.Range(state.PrefixRange(prefix)))

	r, err := state.QueryKeyRanges(ctx, query)
	if err != nil {
		return response.Error(err.Error())
	}
	if r == nil || len(r.([]*state.State)) == 0 {
		return response.NotFound()
	}
	return response.OK(r)
}

func QueryPersonaBasicaHandler(ctx *context.Context) *response.Response {
	args, err := handler.ExtractArgs(ctx.Stub.GetArgs()[1:], CUITParam)
	if err != nil {
		return response.BadRequest("invalid argument: %v", err)
	}
	key := "per:" + strconv.FormatUint(args[0].(uint64), 10) + "#per"
	query := state.Single(state.Point(key))
	r, err := state.QueryKeyRanges(ctx, query)
	if err != nil {
		return response.Error(err.Error())
	}
	if r == nil {
		return response.NotFound()
	}
	return response.OK(r)
}
