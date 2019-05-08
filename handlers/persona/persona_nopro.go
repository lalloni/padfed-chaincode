package persona

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/cast"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func DelPersonaRangeHandler(ctx *context.Context) *response.Response {
	first, err := ctx.ArgUint(1)
	if err != nil {
		return response.BadRequest("invalid persona range start: %v", err)
	}
	last, err := ctx.ArgUint(2)
	if err != nil {
		return response.BadRequest("invalid persona range end: %v", err)
	}
	ids, err := ctx.Store.DelCompositeRange(cast.Persona, store.Range{First: first, Last: last})
	if err != nil {
		return response.Error("deleting persona range [%v,%v]: %v", first, last, err)
	}
	return response.OK(ids)
}

func GetPersonaRangeHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}
