package personas

import (
	"bytes"
	"log"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/helpers"
)

func QueryPersonasByRangeFormated(stub shim.ChaincodeStubInterface, cuitInicio string, cuitFin string, full bool) *fabric.Response {
	cuitInicio = "PER_" + cuitInicio
	cuitFin = "PER_" + cuitFin

	if full || (cuitInicio == cuitFin) {
		cuitFin += "z"
	}
	log.Println("Getting from: " + cuitInicio + " to: " + cuitFin)
	resultsIterator, err := stub.GetStateByRange(cuitInicio, cuitFin)
	if err != nil {
		log.Println(err.Error())
		return fabric.SystemErrorResponse(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := buildResponse(resultsIterator, full)
	if err != nil {
		return fabric.SystemErrorResponse(err.Error())
	}
	log.Println("- query:" + buffer.String())

	return fabric.SuccessResponseWithBuffer(&buffer)
}

func buildResponse(resultsIterator shim.StateQueryIteratorInterface, full bool) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	var currentKey []string

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer, err
		}
		currentKey = (strings.Split(queryResponse.Key, "_"))

		if !full && len(currentKey) > 2 {
			// ignore assets diff to persona
			continue
		}
		helpers.WriteInBuffer(&buffer, queryResponse.Value, queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return buffer, nil
}
