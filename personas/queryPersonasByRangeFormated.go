package personas

import (
	"bytes"
	"log"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/helpers"
)

func QueryPersonasByRangeFormated(stub shim.ChaincodeStubInterface, cuitInicio string, cuitFin string, full bool, composed bool) *fabric.Response {
	cuitInicio = "PER_" + cuitInicio
	cuitFin = "PER_" + cuitFin

	if full || composed || (cuitInicio == cuitFin) {
		cuitFin += "z"
	}
	log.Println("Getting from: " + cuitInicio + " to: " + cuitFin)
	resultsIterator, err := stub.GetStateByRange(cuitInicio, cuitFin)
	if err != nil {
		log.Println(err.Error())
		return fabric.SystemErrorResponse(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := buildResponse(resultsIterator, full, composed)
	if err != nil {
		return fabric.SystemErrorResponse(err.Error())
	}
	log.Println("- query:" + buffer.String())

	return fabric.SuccessResponseWithBuffer(&buffer)
}

func buildResponse(resultsIterator shim.StateQueryIteratorInterface, full bool, composed bool) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	var bufferPersona, bufferImpuesto []byte
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
		if composed {
			if strings.Compare(currentRecord[1], previusCuit) != 0 {
				bufferPersona, bufferImpuesto, bArrayMemberAlreadyWritten = checkBuffers(&buffer, bufferPersona, bufferImpuesto, previusCuit, bArrayMemberAlreadyWritten)
			}
			if len(currentRecord) == 2 {
				bufferPersona = queryResponse.Value
			} else {
				if len(bufferImpuesto) > 0 {
					bufferImpuesto = append(bufferImpuesto, []byte(",")...)
				}
				bufferImpuesto = append(bufferImpuesto, queryResponse.Value...)
			}
			previusCuit = currentRecord[1]
		} else {
			if !full && len(currentRecord) > 2 {
				// ignore assets diff to persona
				continue
			}
			helpers.WriteInBuffer(&buffer, queryResponse.Value, queryResponse.Key, bArrayMemberAlreadyWritten)
			bArrayMemberAlreadyWritten = true
		}
	}
	_, _, _ = checkBuffers(&buffer, bufferPersona, bufferImpuesto, previusCuit, bArrayMemberAlreadyWritten)
	buffer.WriteString("]")
	return buffer, nil
}

var (
	cierre  = []byte("}")
	imps    = []byte(",\"impuestos\":[")
	cierres = []byte("]}")
)

func checkBuffers(buffer *bytes.Buffer, bufferPersona []byte, bufferImpuesto []byte, previusCuit string, bArrayMemberAlreadyWritten bool) ([]byte, []byte, bool) {
	if len(bufferPersona) > 0 {
		if len(bufferImpuesto) > 0 {
			bufferPersona = bytes.Replace(bufferPersona, cierre, imps, 1)
			bufferPersona = append(append(bufferPersona, bufferImpuesto...), cierres...)
		}
		helpers.WriteInBuffer(buffer, bufferPersona, "PER_"+previusCuit, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	return bufferPersona, bufferImpuesto, bArrayMemberAlreadyWritten
}
