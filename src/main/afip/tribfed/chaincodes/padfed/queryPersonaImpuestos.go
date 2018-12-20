package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryPersonaImpuestos(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Numero incorrecto de parametros. Se espera {CUIT}")
	}
	_, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return shim.Error("CUIT [" + args[0] + "] invalido. " + err.Error())
	}
	return s.queryByKeyRange(APIstub, []string{"PER_" + args[0] + "_IMP_", "PER_" + args[0] + "_IMP_z"})
}
