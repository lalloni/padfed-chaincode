package personas

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"

	persona "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
)

func QueryPersona(stub shim.ChaincodeStubInterface, _ []string) *fabric.Response {
	var (
		cuit uint64
		full bool
	)
	args := stub.GetArgs()[1:]
	switch len(args) {
	case 1:
		_, err := handler.ExtractArgs(args, persona.CUITParamVar(&cuit))
		if err != nil {
			return fabric.ClientErrorResponse(err.Error())
		}
	case 2:
		_, err := handler.ExtractArgs(args, persona.CUITParamVar(&cuit), param.BoolVar(&full))
		if err != nil {
			return fabric.ClientErrorResponse(err.Error())
		}
	default:
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT>, [P_FULL]}")
	}
	scuit := strconv.FormatUint(cuit, 10)
	return QueryPersonasByRangeFormated(stub, scuit, scuit, full)
}
