package personas

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
)

func DelPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 1 parámetro con {CUIL}")
	}
	// al cuit final se le añade 'z' para barrer con todos los sufijos y completar el rango
	return fabric.DeleteByKeyRange(stub, "PER_"+args[0], "PER_"+args[0]+"z")
}
