package fabric

import (
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
)

func QueryByKeyArgs(stub shim.ChaincodeStubInterface, args []string) *Response {
	return QueryByKey(stub, args[0])
}

func QueryByKey(stub shim.ChaincodeStubInterface, key string) *Response {
	registerAsBytes, err := stub.GetState(key)
	if err != nil {
		return SystemErrorResponse(err.Error())
	}
	r := &Response{}
	if registerAsBytes == nil {
		log.Println("queryByKey:[]")
		r.Buffer.WriteString("[]")
		return r
	}
	r.Buffer.WriteString("[")
	helpers.WriteInBuffer(&r.Buffer, registerAsBytes, key, false)
	r.Buffer.WriteString("]")
	//	log.Println("queryByKey: [" + r.Buffer.String() + "]")
	return r
}
