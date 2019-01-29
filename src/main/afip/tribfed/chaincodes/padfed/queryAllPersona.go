package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryAllPersona(APIstub shim.ChaincodeStubInterface) Response {
	return s.queryPersonasByRangeFormated(APIstub, "", "", false, false)
}
