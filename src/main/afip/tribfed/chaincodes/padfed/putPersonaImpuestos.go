package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) putPersonaImpuestos(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		shim.Error("Numero incorrecto de parametros. Se espera 2 {CUIT, JSON}")
	}

	cuitStr := args[0]

	log.Print("args[0] [" + cuitStr + "]")
	if _, err := getCUITArgs(args); err != nil {
		return shim.Error("CUIT [" + cuitStr + "] invalida")
	}
	if exists, err := keyExists(APIstub, cuitStr); err != nil {
		return shim.Error(err.Error())
	} else if !exists {
		return shim.Error("CUIT [" + cuitStr + "] inexistente")
	}

	var impuestos Impuestos
	if err := json.Unmarshal([]byte(args[1]), &impuestos); err != nil {
		log.Print(args[1])
		return shim.Error("JSON invalido: " + err.Error())
	}

	if rows, err := putPersonaImpuestos(APIstub, cuitStr, impuestos.Impuestos); err != nil {
		log.Print(err.Error())
		return shim.Error(err.Error())
	} else {
		msg := strconv.Itoa(rows) + " assets processed !!!"
		log.Print(msg)
		return shim.Success([]byte(msg))
	}
}

func putPersonaImpuestos(APIstub shim.ChaincodeStubInterface, cuit string, impuestos []*Impuesto) (int, error) {

	if hid, impuestoDuplicado := hasDuplicatedImpuestos(impuestos); hid {
		return 0, errors.New("Array con idImpuesto [" + strconv.Itoa(int(impuestoDuplicado.IDImpuesto)) + "] duplicado")
	}

	for _, imp := range impuestos {
		impuestoAsBytes, _ := json.Marshal(imp)

		if !(imp.IDOrganismo == 0 || imp.IDOrganismo == 1 || (imp.IDOrganismo >= 900 && imp.IDOrganismo <= 999)) {
			return 0, errors.New("idOrg [" + strconv.Itoa(int(imp.IDOrganismo)) + "] must be an integer 1:AFIP or between 900 and 999")
		}
		if err := validateIdImpuesto(imp.IDImpuesto); err != nil {
			return 0, err
		}
		if exists, err := existsIdImpuesto(APIstub, imp.IDImpuesto); err != nil {
			return 0, err
		} else if !exists {
			return 0, errors.New("idImpuesto [" + strconv.Itoa(int(imp.IDImpuesto)) + "] no definido en ParamImpuesto")
		}
		if err := validateDate(imp.FechaInscripcion); err != nil {
			return 0, errors.New("fechaInscripcion [" + imp.FechaInscripcion + "]: " + err.Error())
		}
		if imp.Periodo < 190000 || imp.Periodo > 205012 {
			return 0, errors.New("periodo [" + strconv.Itoa(int(imp.Periodo)) + "] must be an integer between 190000 and 205012")
		}

		key := "PER_" + cuit + "_IMP_" + strconv.Itoa(int(imp.IDImpuesto))
		if err := APIstub.PutState(key, impuestoAsBytes); err != nil {
			return 0, errors.New("Error putting key [" + key + "]: " + err.Error())
		}
	}
	return len(impuestos), nil
}
