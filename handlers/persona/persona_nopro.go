package persona

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/meta"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

var DelPersonaRangeHandler = handler.MustFunc(delPersonaRange, param.Uint64, param.Uint64)

func delPersonaRange(ctx *context.Context, first, last uint64) *response.Response {
	ids, err := ctx.Store.DelCompositeRange(meta.Persona, store.R(first, last))
	if err != nil {
		return response.Error("deleting persona range [%v,%v]: %v", first, last, err)
	}
	return response.OK(ids)
}

var GetPersonaRangeHandler = handler.MustFunc(getPersonaRange, param.Uint64, param.Uint64)

func getPersonaRange(ctx *context.Context, first, last uint64) *response.Response {
	ps, err := ctx.Store.GetCompositeRange(meta.Persona, store.R(first, last))
	if err != nil {
		return response.Error("getting persona range [%v,%v]: %v", first, last, err)
	}
	return response.OK(ps)
}

func GetPersonaAllHandler(ctx *context.Context) *response.Response {
	ps, err := ctx.Store.GetCompositeAll(meta.Persona)
	if err != nil {
		return response.Error("getting persona all: %v", err)
	}
	return response.OK(ps)
}
