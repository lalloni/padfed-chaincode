package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryPersona(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	_, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return shim.Error("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	p_full := false
	p_composed := false
	switch len(args) {
	case 1:
		break
	case 2:
		p_full, err = strconv.ParseBool(args[1])
		if err != nil {
			return shim.Error("p_full [" + args[1] + "] invalido. " + err.Error())
		}
	case 3:
		p_full, err = strconv.ParseBool(args[1])
		if err != nil {
			return shim.Error("p_full [" + args[1] + "] invalido. " + err.Error())
		}
		p_composed, err = strconv.ParseBool(args[2])
		if err != nil {
			return shim.Error("p_composed [" + args[2] + "] invalido. " + err.Error())
		}
	default:
		return shim.Error("Número incorrecto de parámetros. Se esperaba {<CUIT>, [P_FULL], [P_COMPOSED]}")
	}
	return s.queryPersonasByRangeFormated(APIstub, args[0], args[0], p_full, p_composed)
}
