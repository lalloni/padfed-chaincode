package personas

import (
	"bytes"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/store/key"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/helpers"
)

func QueryPersonasByRangeFormated(stub shim.ChaincodeStubInterface, cuitInicio string, cuitFin string, full bool) *fabric.Response {
	cuitInicio = "per:" + cuitInicio
	cuitFin = "per:" + cuitFin

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

	return fabric.SuccessResponseWithBuffer(buffer)
}

func buildResponse(resultsIterator shim.StateQueryIteratorInterface, full bool) (*bytes.Buffer, error) {
	buffer := &bytes.Buffer{}
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer, err
		}
		if !full {
			k, err := key.Parse(queryResponse.Key)
			if err != nil {
				return nil, errors.Wrap(err, "parsing asset key")
			}
			if k.Tag.Name != "per" {
				continue
			}
		}
		helpers.WriteInBuffer(buffer, queryResponse.Value, queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return buffer, nil
}
