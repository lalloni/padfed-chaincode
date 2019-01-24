package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) deleteByKeyRange(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 2 {
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba 2 parámetros con {RANGO_INICIO, RANGO_FIN}")
	}
	resultsIterator, err := APIstub.GetStateByRange(args[0], args[1])
	var total int
	var firstDeleted, lastDeleted, partialMsj string
	total = 100
	log.Println("Se eliminaran los primeros 100 elementos para evitar un timeout")
	if err != nil {
		return systemErrorResponse(err.Error())
	}
	var count int
	count = 0
	for resultsIterator.HasNext() {
		count++
		if count > total {
			break
		}
		result, err := resultsIterator.Next()
		if err != nil {
			return systemErrorResponse(err.Error(), count)
		}
		if firstDeleted != "" {
			lastDeleted = result.Key
		} else {
			firstDeleted = result.Key
			if !resultsIterator.HasNext() {
				lastDeleted = firstDeleted
			}
		}
		// lenKeyPER = "PER_" + CUIT = 15 digits
		const lenKeyPER = 15
		if len(result.Key) == lenKeyPER {
			// delete asset CUIT
			regex := *regexp.MustCompile(`(^PER_)(\d{11})`)
			res := regex.MatchString(result.Key)
			if res {
				splitedKey := strings.Split(result.Key, "_")
				APIstub.DelState(splitedKey[1])
			}
		}
		log.Print("[" + strconv.Itoa(count) + "] key to delete [" + result.Key)
		APIstub.DelState(result.Key)
	}
	if firstDeleted != "" {
		partialMsj = "Keys eliminadas desde [" + firstDeleted + "] hasta [" + lastDeleted + "]"
	}
	return successResponse(partialMsj+" hasNext ["+strconv.FormatBool(count > total)+"]", count)
}
