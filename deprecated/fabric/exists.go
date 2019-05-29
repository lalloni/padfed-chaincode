package fabric

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// KeyExists returns true if the key exists
func KeyExists(stub shim.ChaincodeStubInterface, key string) (bool, *Response) {
	Log.Info("Key[" + key + "] using GetState...")
	assetAsByte, err := stub.GetState(key)
	if err != nil {
		return false, SystemErrorResponse(err.Error())
	}
	return assetAsByte != nil, &Response{}
}
