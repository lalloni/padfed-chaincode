package personas

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handlerutil/crud"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
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
