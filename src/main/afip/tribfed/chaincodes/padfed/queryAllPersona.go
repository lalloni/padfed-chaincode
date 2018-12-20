package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) queryAllPersona(APIstub shim.ChaincodeStubInterface) peer.Response {
	return s.queryPersonasByRangeFormated(APIstub, "", "", false, false)
}
