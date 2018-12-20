package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
	"github.com/spacemonkeygo/errors"
)

type ErrorResponse struct {
	Message      string `json:",omitempty"`
	Error        error  `json:",omitempty"`
	ErrorMessage string `json:",omitempty"`
	ErrorStack   string `json:",omitempty"`
}

func (s *SmartContract) annotatedErrorResponse(err error, format string, args ...interface{}) peer.Response {
	r := ErrorResponse{Message: fmt.Sprintf(format, args...)}
	if s.debug {
		r.Error = err
		r.ErrorMessage = err.Error()
		r.ErrorStack = errors.GetStack(err)
	}
	responseAsBytes, _ := json.Marshal(r)
	return shim.Error(string(responseAsBytes))
}

func (s *SmartContract) systemErrorResponse(err error) peer.Response {
	return s.annotatedErrorResponse(err, "Error de sistema")
}

func (s *SmartContract) clientErrorResponse(err error) peer.Response {
	return s.annotatedErrorResponse(err, "Error de cliente/request")
}

func (s *SmartContract) businessErrorResponse(errorStr string) peer.Response {
	return s.annotatedErrorResponse(errors.New(errorStr), "Error de negocio")
}
