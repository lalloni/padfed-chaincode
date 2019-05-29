package persona

import (
	"encoding/json"
	"fmt"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/meta"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

func GetPersonaHandler(ctx *context.Context) *response.Response {

	cuit, err := ctx.ArgUint64(1)
	if err != nil {
		return response.BadRequest("invalid persona id: %v", err)
	}

	per, err := ctx.Store.GetComposite(meta.Persona, cuit)
	if err != nil {
		return response.Error("getting persona: %v", err)
	}
	if per == nil {
		return response.NotFound()
	}

	return response.OK(per)

}

func DelPersonaHandler(ctx *context.Context) *response.Response {

	id, err := ctx.ArgUint64(1)
	if err != nil {
		return response.BadRequest("invalid persona id: %v", err)
	}

	exist, err := ctx.Store.HasComposite(meta.Persona, id)
	if err != nil {
		return response.Error("checking persona existence: %v", err)
	}
	if !exist {
		return response.NotFound()
	}

	err = ctx.Store.DelComposite(meta.Persona, id)
	if err != nil {
		return response.Error("deleting persona: %v", err)
	}

	return response.OK(nil)

}

func PutPersonaHandler(ctx *context.Context) *response.Response {

	bs, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("invalid persona: %v", err)
	}

	res, err := validator.ValidatePersonaJSON(bs)
	if err != nil {
		return response.Error("validating persona: %v", err)
	}
	if !res.Valid() {
		return response.BadRequestWithFault(res)
	}

	per := &model.Persona{}
	err = json.Unmarshal(bs, per)
	if err != nil {
		return response.Error("unmarshalling persona: %v", err)
	}

	return save(ctx, per)

}

func PutPersonaListHandler(ctx *context.Context) *response.Response {

	bs, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("invalid persona list: %v", err)
	}

	res, err := validator.ValidatePersonaListJSON(bs)
	if err != nil {
		return response.Error("validating persona list: %v", err)
	}
	if !res.Valid() {
		return response.BadRequestWithFault(res)
	}

	pers := []model.Persona{}
	err = json.Unmarshal(bs, &pers)
	if err != nil {
		return response.Error("unmarshalling persona list: %v", err)
	}

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
