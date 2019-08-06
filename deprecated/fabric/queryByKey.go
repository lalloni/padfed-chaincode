package fabric

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/helpers"
)

func QueryByKey(stub shim.ChaincodeStubInterface, _ []string) *Response {
	var key string

	_, err := handler.ExtractArgs(stub.GetArgs()[1:], param.StringVar(&key))
	if err != nil {
		return ClientErrorResponse(err.Error())
	}

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
