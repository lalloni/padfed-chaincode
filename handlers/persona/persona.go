package persona

import (
	"fmt"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/meta"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

var GetPersonaHandler = handler.MustFunc(getPersona, persona.CUITParam)

func getPersona(ctx *context.Context, id uint64) *response.Response {

	per, err := ctx.Store.GetComposite(meta.Persona, id)
	if err != nil {
		return response.Error("getting persona: %v", err)
	}
	if per == nil {
		return notFoundResponse(id)
	}

	return response.OK(per)

}

var DelPersonaHandler = handler.MustFunc(delPersona, persona.CUITParam)

func delPersona(ctx *context.Context, id uint64) *response.Response {

	exist, err := ctx.Store.HasComposite(meta.Persona, id)
	if err != nil {
		return response.Error("checking persona existence: %v", err)
	}
	if !exist {
		return notFoundResponse(id)
	}

	err = ctx.Store.DelComposite(meta.Persona, id)
	if err != nil {
		return response.Error("deleting persona: %v", err)
	}

	return response.OK(nil)

}

var PutPersonaHandler = handler.MustFunc(putPersona, persona.PersonaParam)

func putPersona(ctx *context.Context, per *model.Persona) *response.Response {

	return save(ctx, per)

}

var PutPersonaListHandler = handler.MustFunc(putPersonaList, persona.PersonaListParam)

func putPersonaList(ctx *context.Context, pers []model.Persona) *response.Response {

	for n, per := range pers {
		per := per
		res := save(ctx, &per)
		if !res.OK() {
			res.Message = fmt.Sprintf("persona %d: %s", n+1, res.Message)
			return res
		}
	}

	return response.OK(len(pers))

}

func save(ctx *context.Context, per *model.Persona) *response.Response {

	if per.ID == 0 {
		return response.BadRequest("id required")
	}
	if per.Persona != nil && per.ID != per.Persona.ID {
		return response.BadRequest("id %q and persona.id %q must be equal", per.ID, per.Persona.ID)
	}

	exist, err := ctx.Store.HasComposite(meta.Persona, per.ID)
	if err != nil {
		return response.Error("checking persona existence: %v", err)
	}
	if !exist && per.Persona == nil {
		return response.BadRequest("persona is required when putting a new instance")
	}

	err = ctx.Store.PutComposite(meta.Persona, per)
	if err != nil {
		return response.Error("putting persona: %v", err)
	}

	return response.OK(nil)

}

func notFoundResponse(id uint64) *response.Response {
	return response.NotFoundWithMessage("persona with id %d not found", id)
}
