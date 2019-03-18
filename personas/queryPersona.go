package personas

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
)

func QueryPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	var err error
	if _, err = helpers.GetCUIT(args[0]); err != nil {
		return fabric.ClientErrorResponse("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	full := false
	composed := false
	// se esta usando fallthrough, importante no cambiar el orden de los case's
	switch len(args) {
	case 1:
		break
	case 3:
		composed, err = strconv.ParseBool(args[2])
		if err != nil {
			return fabric.ClientErrorResponse("composed [" + args[2] + "] invalido. " + err.Error())
		}
		fallthrough
	case 2:
		full, err = strconv.ParseBool(args[1])
		if err != nil {
			return fabric.ClientErrorResponse("full [" + args[1] + "] invalido. " + err.Error())
		}
	default:
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT>, [P_FULL], [P_COMPOSED]}")
	}
	return QueryPersonasByRangeFormated(stub, args[0], args[0], full, composed)
}
