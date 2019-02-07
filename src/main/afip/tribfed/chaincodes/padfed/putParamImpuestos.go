package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) putParamImpuestos(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 1 {
		return clientErrorResponse("Número incorrecto de parámetros. Se espera un parámetro con {[]JSON}")
	}
	var err error

	var impuestos []ParamImpuesto
	if err = json.Unmarshal([]byte(args[0]), &impuestos); err != nil {
		msg := "JSON invalido: " + err.Error()
		log.Println(msg)
		return systemErrorResponse(msg)
	}
	count := 0
	for _, imp := range impuestos {
		if !(imp.IDOrganismo == 1 || (imp.IDOrganismo >= 900 && imp.IDOrganismo <= 999)) {
			return clientErrorResponse("idOrg ["+strconv.Itoa(int(imp.IDOrganismo))+"] debe ser un entero igual a 1:AFIP o entre 900 y 999", count)
		}

		if err := validateIdImpuesto(imp.Impuesto); err.isError() {
			err.WrongItem = count
			return err
		}

		impuestoAsBytes, _ := json.Marshal(imp)
		var key = getParamImpuestoKey(imp.Impuesto)
		if err := APIstub.PutState(key, impuestoAsBytes); err != nil {
			return systemErrorResponse("Error putting key ["+key+"] - "+err.Error(), count)
		}
		count++
	}
	log.Println(strconv.Itoa(len(impuestos)) + " assets processed !!!")
	return successResponse("Ok", len(impuestos))
}

// validateIdImpuesto valida el codigo de impuesto sin acceder al State
func validateIdImpuesto(impuesto int32) Response {
	if !(impuesto >= 1 && impuesto <= 9999) {
		return clientErrorResponse("impuesto [" + strconv.Itoa(int(impuesto)) + "] debe ser un entero entre 1 y 9999")
	}
	return Response{}
}

// existsIdImpuesto verifica que exista un asset "IMP_<idImpuesot>"
func existsIdImpuesto(APIstub shim.ChaincodeStubInterface, impuesto int32) (bool, Response) {
	if exists, err := keyExists(APIstub, getParamImpuestoKey(impuesto)); err.isError() {
		return false, err
	} else {
		return exists, Response{}
	}
}

func getParamImpuestoKey(impuesto int32) string {
	return "IMP_" + strconv.Itoa(int(impuesto))
}
