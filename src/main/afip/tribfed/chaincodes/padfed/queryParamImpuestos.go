package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) queryParamImpuestos(APIstub shim.ChaincodeStubInterface) Response {
	return s.queryByKeyRange(APIstub, []string{"IMP_", "IMP_z"})
}
