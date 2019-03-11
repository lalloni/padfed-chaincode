package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) delPersonaAssets(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 2 {
		return clientErrorResponse("Numero incorrecto de parametros. Se esperaba 2 parametros con {CUIL, []KEYS}")
	}
	var cuit = args[0]
	var keys []string
	if err := json.Unmarshal([]byte(args[1]), &keys); err != nil {
		msg := "JSON de array invalido: " + err.Error()
		log.Println(msg)
		return systemErrorResponse(msg)
	}

	if len(keys) == 0 {
		return clientErrorResponse("El array debe tener por lo menos un elemento")
	}
	if isErr, errResponse := checkDuplicated(keys, cuit); isErr {
		return errResponse
	}
	if exists, err := keyExists(APIstub, "PER_"+cuit); err.isError() {
		return err
	} else {
		if !exists {
			return clientErrorResponse("Debe existir como un asset la persona: " + cuit)
		}
	}
	count := 0
	for _, key := range keys {
		log.Print("key to delete [" + key + "]")
		if err := APIstub.DelState(key); err != nil {
			return systemErrorResponse("Error al eliminar: [" + key + "] " + err.Error())
		}
		count++
	}

	log.Println(strconv.Itoa(count) + " assets deleted !!!")
	return successResponse("Ok", count)
}

func checkDuplicated(array []string, cuit string) (bool, Response) {
	var FORMAT_KEY_REGEXP = *regexp.MustCompile(`^(PER_)` + cuit + `(_IMP_)(\d+)$`)
	keys := make(map[string]bool)
	for _, entry := range array {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res := FORMAT_KEY_REGEXP.FindStringSubmatch(entry)
			if len(res) != 4 {
				return true, clientErrorResponse("El array debe tener keys compuestas de assets PersonaImpuesto PER_" + cuit + "_IMP_[imp]")
			}
		} else {
			return true, clientErrorResponse("El array no puede tener elementos repetidos")
		}
	}
	return false, Response{}
}
