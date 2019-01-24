package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) delParamImpuestosAll(APIstub shim.ChaincodeStubInterface) Response {
	return s.deleteByKeyRange(APIstub, []string{"IMP_", "IMP_z"})
}
