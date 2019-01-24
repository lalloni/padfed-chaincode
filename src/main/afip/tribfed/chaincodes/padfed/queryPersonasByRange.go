package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryPersonasByRange(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	_, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return s.peerResponse(clientErrorResponse("CUIT [" + args[0] + "] invalido. " + err.Error()))
	}
	if len(args[0]) != 11 {
		return s.peerResponse(clientErrorResponse("Longitud del CUIT [" + args[0] + "], invalido."))
	}
	_, err = strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return s.peerResponse(clientErrorResponse("CUIT [" + args[1] + "] invalido. " + err.Error()))
	}
	if len(args[1]) != 11 {
		return s.peerResponse(clientErrorResponse("Longitud del CUIT [" + args[1] + "], invalido."))
	}
	switch len(args) {
	case 2:
		return s.queryPersonasByRangeFormated(APIstub, args[0], args[1], false, false)
	case 3:
		p_full, err1 := strconv.ParseBool(args[2])
		if err1 != nil {
			return s.peerResponse(clientErrorResponse(err1.Error()))
		}
		return s.queryPersonasByRangeFormated(APIstub, args[0], args[1], p_full, false)
	case 4:
		p_full, err1 := strconv.ParseBool(args[2])
		if err1 != nil {
			return s.peerResponse(clientErrorResponse("p_full [" + args[2] + "] invalido. " + err.Error()))
		}
		p_composed, err2 := strconv.ParseBool(args[3])
		if err2 != nil {
			return s.peerResponse(clientErrorResponse("p_composed [" + args[3] + "] invalido. " + err.Error()))
		}
		return s.queryPersonasByRangeFormated(APIstub, args[0], args[1], p_full, p_composed)
	default:
		return s.peerResponse(clientErrorResponse("Número incorrecto de parámetros. Se esperaba {<CUIT_INICIO>, <CUIT_FIN>, [P_FULL], [P_COMPOSED]}"))
	}
}
