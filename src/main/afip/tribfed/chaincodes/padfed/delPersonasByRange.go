package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) delPersonasByRange(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 2 {
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba 2 parámetros con {CUIL_INICIO, CUIL_FIN}")
	}
	// al cuit final se le añade 'z' para barrer con todos los sufijos y completar el rango
	return s.deleteByKeyRange(APIstub, []string{"PER_" + args[0], "PER_" + args[1] + "z"})
}
