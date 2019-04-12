package fabric

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/helpers"
)

func QueryByKeyArgs(stub shim.ChaincodeStubInterface, args []string) *Response {
	return QueryByKey(stub, args[0])
}

func QueryByKey(stub shim.ChaincodeStubInterface, key string) *Response {
	registerAsBytes, err := stub.GetState(key)
	if err != nil {
		return SystemErrorResponse(err.Error())
	}
	var b bytes.Buffer
	if registerAsBytes == nil {
		log.Println("queryByKey:[]")
		b.WriteString("[]")
		return SuccessResponseWithBuffer(&b)
	}
	b.WriteString("[")
	helpers.WriteInBuffer(&b, registerAsBytes, key, false)
	b.WriteString("]")
	return SuccessResponseWithBuffer(&b)
}
