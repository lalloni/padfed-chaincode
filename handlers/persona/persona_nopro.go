package persona

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/cast"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func DelPersonaRangeHandler(ctx *context.Context) *response.Response {
	r, first, last, err := argUintRange(ctx, 1)
	if err != nil {
		return response.BadRequest("invalid persona range: %v", err)
	}
	ids, err := ctx.Store.DelCompositeRange(cast.Persona, r)
	if err != nil {
		return response.Error("deleting persona range [%v,%v]: %v", first, last, err)
	}
	return response.OK(ids)
}

func GetPersonaRangeHandler(ctx *context.Context) *response.Response {
	r, first, last, err := argUintRange(ctx, 1)
	if err != nil {
		return response.BadRequest("invalid persona range: %v", err)
	}
	ps, err := ctx.Store.GetCompositeRange(cast.Persona, r)
	if err != nil {
		return response.Error("getting persona range [%v,%v]: %v", first, last, err)
	}
	return response.OK(ps)
}

func argUintRange(ctx *context.Context, pos int) (*store.Range, uint, uint, error) {
	a, err := ctx.ArgUint(pos)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "invalid range start")
	}
	b, err := ctx.ArgUint(pos + 1)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "invalid range end")
	}
	return &store.Range{First: a, Last: b}, a, b, nil
}
