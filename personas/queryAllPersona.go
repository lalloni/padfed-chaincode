package personas

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
)

func QueryAllPersona(stub shim.ChaincodeStubInterface, _ []string) *fabric.Response {
	return QueryPersonasByRangeFormated(stub, "", "", false)
}
