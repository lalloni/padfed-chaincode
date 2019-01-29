package main

import (
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryByKey(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 1 {
		return clientErrorResponse("Numero incorrecto de parametros. Se espera {KEY}")
	}
	registerAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return systemErrorResponse(err.Error())
	}
	var r Response
	if registerAsBytes == nil {
		log.Println("queryByKey:[]")
		r.Buffer.WriteString("[]")
		return r
	}
	r.Buffer.WriteString("[")
	writeInBuffer(&r.Buffer, string(registerAsBytes), args[0], false)
	r.Buffer.WriteString("]")
	//	log.Println("queryByKey: [" + r.Buffer.String() + "]")
	return r
}
