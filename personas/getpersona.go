package personas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func GetPersonaAPI(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	cuit, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("cuit inválido: %v", args[0]))
	}
	st := store.New(stub)
	p, err := st.GetComposite(Persona, Persona.IdentifierKey(cuit))
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("obteniendo persona: %v", err))
	}
	if p == nil {
		return fabric.NotFoundErrorResponse()
	}
	p.(*model.Persona).ID = p.(*model.Persona).Persona.ID
	bs, err := json.Marshal(p)
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("generando respuesta: %v", err))
	}
	return fabric.SuccessResponseWithBuffer(bytes.NewBuffer(bs))
}
