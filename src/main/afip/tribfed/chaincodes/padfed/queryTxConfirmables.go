package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryTxConfirmables(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 || len(args) > 4 {
		return s.peerResponse(clientErrorResponse("Numero incorrecto de parametros. Se espera {ID_ORG, ID_TXC, CUIT, PENDIENTE}"))
	}
	idOrg, err := strconv.Atoi(args[0])
	if err != nil {
		return s.peerResponse(clientErrorResponse("ID_ORG [" + args[0] + "] invalido. " + err.Error()))
	}
	if idOrg < 900 || idOrg > 999 {
		return s.peerResponse(clientErrorResponse("ID_ORG incorrecto. Se espera el rango [900...999]"))
	}
	idTxc := ""
	cuit := ""
	isPendiente := "SIN_DEFINIR"
	switch len(args) {
	case 2:
		err = checkParam(args[1])
		if err != nil {
			return s.peerResponse(clientErrorResponse("ID_TXC [" + args[1] + "] invalido. " + err.Error()))
		}
		idTxc = args[1]
	case 3:
		err = checkParam(args[1])
		if err != nil {
			return s.peerResponse(clientErrorResponse("ID_TXC [" + args[1] + "] invalido. " + err.Error()))
		}
		err = checkParam(args[2])
		if err != nil {
			return s.peerResponse(clientErrorResponse("CUIT [" + args[2] + "] invalido. " + err.Error()))
		}
		idTxc = args[1]
		cuit = args[2]
	case 4:
		err = checkParam(args[1])
		if err != nil {
			return s.peerResponse(clientErrorResponse("ID_TXC [" + args[1] + "] invalido. " + err.Error()))
		}
		err = checkParam(args[2])
		if err != nil {
			return s.peerResponse(clientErrorResponse("CUIT [" + args[2] + "] invalido. " + err.Error()))
		}
		isPendiente = args[3]
		idTxc = args[1]
		cuit = args[2]
	}
	startKey := "ORG_" + args[0] + "_TXC_"
	if idTxc != "" {
		startKey += idTxc
	}
	endKey := startKey + "z"
	return s.queryByKeyRangeWithFilters(APIstub, startKey, endKey, cuit, isPendiente)
}

func checkParam(param string) error {
	if param != "" {
		_, err := strconv.ParseInt(param, 10, 64)
		return err
	}
	return nil
}

func (s *SmartContract) queryByKeyRangeWithFilters(APIstub shim.ChaincodeStubInterface, startKey string, endKey string, cuit string, isPendiente string) peer.Response {
	log.Println("Getting from: " + startKey + " to: " + endKey)
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		log.Println(err.Error())
		return s.peerResponse(systemErrorResponse(err.Error()))
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return s.peerResponse(systemErrorResponse(err.Error()))
		}
		var txConfirmable TXConfirmable

		err = json.Unmarshal(queryResponse.Value, &txConfirmable)
		if err != nil {
			return s.peerResponse(systemErrorResponse("JSON invalido: " + err.Error()))
		}
		if cuit != "" {
			if strings.Compare(strconv.FormatUint(txConfirmable.CUIT, 10), cuit) != 0 {
				continue
			}
		}
		if isPendiente == "true" {
			if txConfirmable.TipoRespuesta == 1 || txConfirmable.TipoRespuesta == 2 {
				continue
			}
		} else if isPendiente == "false" {
			if txConfirmable.TipoRespuesta != 1 && txConfirmable.TipoRespuesta != 2 {
				continue
			}
		}
		writeInBuffer(&buffer, string(queryResponse.Value), queryResponse.Key, bArrayMemberAlreadyWritten)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	log.Println("- queryTxConfirmables:" + buffer.String())
	return shim.Success(buffer.Bytes())
}
