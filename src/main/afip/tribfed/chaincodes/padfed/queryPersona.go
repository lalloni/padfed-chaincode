package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryPersona(APIstub shim.ChaincodeStubInterface, args []string) Response {
	var err error
	if _, err = getCUITArgs(args); err != nil {
		return clientErrorResponse("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	p_full := false
	p_composed := false
	switch len(args) {
	case 1:
		break
	case 2:
		p_full, err = strconv.ParseBool(args[1])
		if err != nil {
			return clientErrorResponse("p_full [" + args[1] + "] invalido. " + err.Error())
		}
	case 3:
		p_full, err = strconv.ParseBool(args[1])
		if err != nil {
			return clientErrorResponse("p_full [" + args[1] + "] invalido. " + err.Error())
		}
		p_composed, err = strconv.ParseBool(args[2])
		if err != nil {
			return clientErrorResponse("p_composed [" + args[2] + "] invalido. " + err.Error())
		}
	default:
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT>, [P_FULL], [P_COMPOSED]}")
	}
	return s.queryPersonasByRangeFormated(APIstub, args[0], args[0], p_full, p_composed)
}
