package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) putParamImpuestos(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Número incorrecto de parámetros. Se espera un parámetro con {[]JSON}")
	}
	var err error

	var impuestos []ParamImpuesto
	if err = json.Unmarshal([]byte(args[0]), &impuestos); err != nil {
		msg := "JSON invalido: " + err.Error()
		log.Println(msg)
		return shim.Error(msg)
	}
	for _, imp := range impuestos {

		if !(imp.IDOrganismo == 1 || (imp.IDOrganismo >= 900 && imp.IDOrganismo <= 999)) {
			return s.businessErrorResponse("idOrg [" + strconv.Itoa(int(imp.IDOrganismo)) + "] debe ser un entero igual a 1:AFIP o entre 900 y 999")
		}

		if err := validateIdImpuesto(imp.IDImpuesto); err != nil {
			return s.businessErrorResponse(err.Error())
		}

		impuestoAsBytes, _ := json.Marshal(imp)
		var key = getParamImpuestoKey(imp.IDImpuesto)
		if err := APIstub.PutState(key, impuestoAsBytes); err != nil {
			return shim.Error("Error putting key [" + key + "] - " + err.Error())
		}
	}
	msg := strconv.Itoa(len(impuestos)) + " assets processed !!!"
	log.Println(msg)
	return shim.Success([]byte(msg))
}

// validateIdImpuesto valida el codigo de impuesto sin acceder al State
func validateIdImpuesto(idImpuesto int32) error {
	if !(idImpuesto >= 1 && idImpuesto <= 9999) {
		return errors.New("idImpuesto [" + strconv.Itoa(int(idImpuesto)) + "] debe ser un entero entre 1 y 9999")
	}
	return nil
}

// existsIdImpuesto verifica que exista un asset "IMP_<idImpuesot>"
func existsIdImpuesto(APIstub shim.ChaincodeStubInterface, idImpuesto int32) (bool, error) {
	if exists, err := keyExists(APIstub, getParamImpuestoKey(idImpuesto)); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func getParamImpuestoKey(idImpuesto int32) string {
	return "IMP_" + strconv.Itoa(int(idImpuesto))
}
