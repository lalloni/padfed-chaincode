package persona

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/persona"
	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestingHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(
		common.Defaults, // i.e. get, getrange, has, put, putlist, del, delrange
		common.WithIDParam(CUITParam),
		common.WithItemParam(PersonaParam),
		common.WithListParam(PersonaListParam),
		common.WithValidator(validatePersona),
	)
	if !testing {
		opts = append(opts, common.WithWriteCheck(common.AFIP))
	}
	common.AddCRUDHandlers(r, model.Schema, opts...)
}

func validatePersona(ctx *context.Context, v interface{}) *response.Response {

	per := v.(*model.Persona)

	if per.ID == 0 {
		return response.BadRequest("id required")
	}

	if per.Persona != nil && per.ID != per.Persona.ID {
		return response.BadRequest("id %q and persona.id %q must be equal", per.ID, per.Persona.ID)
	}

	exist, err := ctx.Store.HasComposite(persona.Schema, per.ID)
	if err != nil {
		return response.Error("checking persona existence: %v", err)
	}

	if !exist && per.Persona == nil {
		return response.BadRequest("persona is required when putting a new instance")
	}

	return nil

}
