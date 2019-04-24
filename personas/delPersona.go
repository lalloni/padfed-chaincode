package personas

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/cast"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func DelPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 1 parámetro con {CUIL}")
	}
	cuit, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("cuit inválido: %v", args[0]))
	}

	st := store.New(stub)

	if exist, err := st.HasComposite(cast.Persona, cuit); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Error obteniendo existencia de persona: %v", err))
	} else if !exist {
		return fabric.NotFoundErrorResponse()
	}

	err = st.DelComposite(cast.Persona, cuit)
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("eliminando persona: %v", err))
	}
	return fabric.SuccessResponse("persona eliminada", 1)
}
