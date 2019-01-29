package main

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryByKeyRange(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 2 {
		return clientErrorResponse("Numero incorrecto de parametros. Se espera {START_KEY, END_KEY}")
	}
	START_KEY := args[0]
	END_KEY := args[1] + "z"

	log.Println("Getting from: " + START_KEY + " to: " + END_KEY)
	resultsIterator, err := APIstub.GetStateByRange(START_KEY, END_KEY)
	if err != nil {
		log.Println(err.Error())
		return systemErrorResponse(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return systemErrorResponse(err.Error())
		}
		writeInBuffer(&buffer, string(queryResponse.Value), queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	log.Println("- queryByKeyRange:" + buffer.String())
	return successResponseWithBuffer(&buffer)
}
