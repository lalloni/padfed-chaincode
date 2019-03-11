package inscripciones

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
)

func QueryPersonaImpuestos(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Numero incorrecto de parametros. Se espera {CUIT}")
	}
	if _, err := helpers.GetCUIT(args[0]); err != nil {
		return fabric.ClientErrorResponse("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	return fabric.QueryByKeyRange(stub, "PER_"+args[0]+"_IMP_", "PER_"+args[0]+"_IMP_z")
}
