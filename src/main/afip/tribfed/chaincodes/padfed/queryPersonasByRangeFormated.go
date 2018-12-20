package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryPersonasByRangeFormated(APIstub shim.ChaincodeStubInterface, cuit_inicio string, cuit_fin string, p_full bool, p_composed bool) peer.Response {
	cuit_inicio = "PER_" + cuit_inicio
	cuit_fin = "PER_" + cuit_fin

	if p_full || p_composed || (cuit_inicio == cuit_fin) {
		cuit_fin += "z"
	}
	log.Println("Getting from: " + cuit_inicio + " to: " + cuit_fin)
	resultsIterator, err := APIstub.GetStateByRange(cuit_inicio, cuit_fin)
	if err != nil {
		log.Println(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := buildResponse(resultsIterator, p_full, p_composed)
	if err != nil {
		return shim.Error(err.Error())
	}
	log.Println("- query:" + buffer.String())

	return shim.Success(buffer.Bytes())
}

func buildResponse(resultsIterator shim.StateQueryIteratorInterface, p_full bool, p_composed bool) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	var bufferPersona, bufferImpuesto string
	bArrayMemberAlreadyWritten := false
	previusCuit := ""
	var currentRecord []string

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer, err
		}
		currentRecord = (strings.Split(queryResponse.Key, "_"))
		if p_composed == true {
			if strings.Compare(currentRecord[1], previusCuit) != 0 {
				bufferPersona, bufferImpuesto, bArrayMemberAlreadyWritten = checkBuffers(&buffer, bufferPersona, bufferImpuesto, previusCuit, bArrayMemberAlreadyWritten)
			}
			if len(currentRecord) == 2 {
				bufferPersona = string(queryResponse.Value)
			} else {
				if len(bufferImpuesto) > 0 {
					bufferImpuesto += ","
				}
				bufferImpuesto += string(queryResponse.Value)
			}
			previusCuit = currentRecord[1]
		} else {
			if p_full == false && len(currentRecord) > 2 {
				// ignore assets diff to persona
				continue
			}
			writeInBuffer(&buffer, string(queryResponse.Value), queryResponse.Key, bArrayMemberAlreadyWritten)
			bArrayMemberAlreadyWritten = true
		}
	}
	_, _, _ = checkBuffers(&buffer, bufferPersona, bufferImpuesto, previusCuit, bArrayMemberAlreadyWritten)
	buffer.WriteString("]")
	return buffer, nil
}

func checkBuffers(buffer *bytes.Buffer, bufferPersona string, bufferImpuesto string, previusCuit string, bArrayMemberAlreadyWritten bool) (string, string, bool) {
	if len(bufferPersona) > 0 {
		if len(bufferImpuesto) > 0 {
			bufferPersona = strings.Replace(bufferPersona, "}", ",\"impuestos\":[", 1)
			bufferPersona += bufferImpuesto + "]}"
		}
		writeInBuffer(buffer, bufferPersona, "PER_"+previusCuit, bArrayMemberAlreadyWritten)
		// clean buffers
		bufferPersona = ""
		bufferImpuesto = ""
		bArrayMemberAlreadyWritten = true
	}
	return bufferPersona, bufferImpuesto, bArrayMemberAlreadyWritten
}
