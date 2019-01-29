package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryPersonaImpuestos(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 1 {
		return clientErrorResponse("Numero incorrecto de parametros. Se espera {CUIT}")
	}
	if _, err := getCUITArgs(args); err != nil {
		return clientErrorResponse("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	return s.queryByKeyRange(APIstub, []string{"PER_" + args[0] + "_IMP_", "PER_" + args[0] + "_IMP_z"})
}
