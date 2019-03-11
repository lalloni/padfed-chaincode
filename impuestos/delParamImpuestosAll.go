package impuestos

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
)

func DeleteAll(stub shim.ChaincodeStubInterface) *fabric.Response {
	return fabric.DeleteByKeyRange(stub, "IMP_", "IMP_z")
}
