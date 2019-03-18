package personas

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
)

func QueryPersonasByRange(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	var err error
	for _, cuit := range []string{args[0], args[1]} {
		_, err = strconv.ParseUint(cuit, 10, 64)
		if err != nil {
			return fabric.ClientErrorResponse("CUIT [" + cuit + "] invalido. " + err.Error())
		}
		if len(cuit) != 11 {
			return fabric.ClientErrorResponse("Longitud del CUIT [" + cuit + "], invalido.")
		}
	}
	full := false
	composed := false
	// se esta usando fallthrough, importante no cambiar el orden de los case's
	switch len(args) {
	case 4:
		composed, err = strconv.ParseBool(args[3])
		if err != nil {
			return fabric.ClientErrorResponse("composed [" + args[3] + "] invalido. " + err.Error())
		}
		fallthrough
	case 3:
		full, err = strconv.ParseBool(args[2])
		if err != nil {
			return fabric.ClientErrorResponse("full [" + args[2] + "] invalido. " + err.Error())
		}
		fallthrough
	case 2:
		return QueryPersonasByRangeFormated(stub, args[0], args[1], full, composed)
	default:
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT_INICIO>, <CUIT_FIN>, [P_FULL], [P_COMPOSED]}")
	}
}
