package fabric

import (
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/helpers"
)

func QueryByKeyRange(stub shim.ChaincodeStubInterface, _ []string) *Response {

	args, err := handler.ExtractArgs(stub.GetArgs()[1:], param.String, param.String)
	if err != nil {
		return ClientErrorResponse(err.Error())
	}

	startKey := args[0].(string)
	endKey := args[1].(string) + "z"

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
