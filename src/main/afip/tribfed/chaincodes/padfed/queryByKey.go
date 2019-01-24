package main

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryByKey(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return s.peerResponse(clientErrorResponse("Numero incorrecto de parametros. Se espera {KEY}"))
	}
	registerAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return s.peerResponse(systemErrorResponse(err.Error()))
	} else if registerAsBytes == nil {
		log.Println("- queryByKey:[]")
		return shim.Success([]byte("[]"))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	writeInBuffer(&buffer, string(registerAsBytes), args[0], false)
	buffer.WriteString("]")
	log.Println("- queryByKey: [" + buffer.String() + "]")
	return shim.Success(buffer.Bytes())
}
