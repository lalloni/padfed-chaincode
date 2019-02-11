package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryPersonasByRange(APIstub shim.ChaincodeStubInterface, args []string) Response {
	var err error
	for _, cuit := range []string{args[0], args[1]} {
		_, err = strconv.ParseUint(cuit, 10, 64)
		if err != nil {
			return clientErrorResponse("CUIT [" + cuit + "] invalido. " + err.Error())
		}
		if len(cuit) != 11 {
			return clientErrorResponse("Longitud del CUIT [" + cuit + "], invalido.")
		}
	}
	p_full := false
	p_composed := false
	// se esta usando fallthrough, importante no cambiar el orden de los case's
	switch len(args) {
	case 4:
		p_composed, err = strconv.ParseBool(args[3])
		if err != nil {
			return clientErrorResponse("p_composed [" + args[3] + "] invalido. " + err.Error())
		}
		fallthrough
	case 3:
		p_full, err = strconv.ParseBool(args[2])
		if err != nil {
			return clientErrorResponse("p_full [" + args[2] + "] invalido. " + err.Error())
		}
		fallthrough
	case 2:
		return s.queryPersonasByRangeFormated(APIstub, args[0], args[1], p_full, p_composed)
	default:
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT_INICIO>, <CUIT_FIN>, [P_FULL], [P_COMPOSED]}")
	}
}
