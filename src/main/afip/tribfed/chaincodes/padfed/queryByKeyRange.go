package main

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryByKeyRange(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Numero incorrecto de parametros. Se espera {START_KEY, END_KEY}")
	}
	START_KEY := args[0]
	END_KEY := args[1] + "z"

	log.Println("Getting from: " + START_KEY + " to: " + END_KEY)
	resultsIterator, err := APIstub.GetStateByRange(START_KEY, END_KEY)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		writeInBuffer(&buffer, string(queryResponse.Value), queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	log.Println("- queryByKeyRange:" + buffer.String())
	return shim.Success(buffer.Bytes())
}
