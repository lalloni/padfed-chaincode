package personas

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	persona "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
)

func DelPersonasByRange(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 2 parámetros con {CUIL_INICIO, CUIL_FIN}")
	}
	sid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Parámetro 1 incorrecto: %q: %v", args[0], err))
	}
	eid, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Parámetro 2 incorrecto: %q: %v", args[1], err))
	}
	sk, err := persona.Schema.IdentifierKey(sid)
	if err != nil {
		return fabric.SystemErrorResponse(err.Error())
	}
	ek, err := persona.Schema.IdentifierKey(eid)
	if err != nil {
		return fabric.SystemErrorResponse(err.Error())
	}
	start, _ := sk.Range()
	_, end := ek.Range()
	return fabric.DeleteByKeyRange(stub, []string{start, end})
}
