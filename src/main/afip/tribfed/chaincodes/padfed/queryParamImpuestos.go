package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryParamImpuestos(APIstub shim.ChaincodeStubInterface) peer.Response {
	return s.queryByKeyRange(APIstub, []string{"IMP_", "IMP_z"})
}
