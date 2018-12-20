package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) delPersonasByRange(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Número incorrecto de parámetros. Se esperaba 2 parámetros con {CUIL_INICIO, CUIL_FIN}")
	}
	// al cuit final se le añade 'z' para barrer con todos los sufijos y completar el rango
	return s.deleteByKeyRange(APIstub, []string{"PER_" + args[0], "PER_" + args[1] + "z"})
}
