package personas

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
)

func DelPersonasByRange(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 2 parámetros con {CUIL_INICIO, CUIL_FIN}")
	}
	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Parámetro 1 incorrecto: %q: %v", args[0], err))
	}
	start, _ := Persona.IdentifierKey(id).Range()
	id, err = strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Parámetro 2 incorrecto: %q: %v", args[1], err))
	}
	_, end := Persona.IdentifierKey(id).Range()
	return fabric.DeleteByKeyRange(stub, start, end)
}
