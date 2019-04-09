package personas

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
)

func DelPersonasByRange(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 2 parámetros con {CUIL_INICIO, CUIL_FIN}")
	}
	// al cuit final se le añade 'z' para barrer con todos los sufijos y completar el rango
	return fabric.DeleteByKeyRange(stub, "PER_"+args[0], "PER_"+args[1]+"z")
}
