package fabric

import (
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/helpers"
)

func QueryByKeyRangeArgs(stub shim.ChaincodeStubInterface, args []string) *Response {
	return QueryByKeyRange(stub, args[0], args[1])
}

func QueryByKeyRange(stub shim.ChaincodeStubInterface, startKey, endKey string) *Response {
	endKey += "z" // TODO esto parece un bug: el llamador debería decidir la end key!

	Log.Info("Getting from: " + startKey + " to: " + endKey)
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return SystemErrorResponse(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return SystemErrorResponse(err.Error())
		}
		helpers.WriteInBuffer(&buffer, queryResponse.Value, queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	Log.Info("- queryByKeyRange:" + buffer.String())
	return SuccessResponseWithBuffer(&buffer)
}
