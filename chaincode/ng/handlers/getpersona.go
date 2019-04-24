package handlers

import (
	"encoding/json"

	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/support"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/personas"
)

func GetPersonaHandler(ctx *support.Context) peer.Response {
	cuit, err := ctx.ArgUint(1)
	if err != nil {
		return response.BadRequest("invalid persona id")
	}
	p, err := ctx.Store.GetComposite(personas.Persona, cuit)
	if err != nil {
		return response.ErrorWrap(err, "getting persona")
	}
	if p == nil {
		return response.NotFound("persona with id %v", cuit)
	}
	bs, err := json.Marshal(p)
	if err != nil {
		return response.ErrorWrap(err, "marshaling persona")
	}
	return response.Success(bs)
}
