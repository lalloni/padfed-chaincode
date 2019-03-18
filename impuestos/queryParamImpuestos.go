package impuestos

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
)

func QueryParamImpuestos(stub shim.ChaincodeStubInterface, _ []string) *fabric.Response {
	return fabric.QueryByKeyRange(stub, "IMP_", "IMP_z")
}
