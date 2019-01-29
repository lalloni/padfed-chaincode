package main

import (
	"bytes"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

const (
	// code response
	OK                    int32 = 200
	BAD_REQUEST           int32 = 400
	FORBIDDEN             int32 = 403
	INTERNAL_SERVER_ERROR int32 = 500
)

type Response struct {
	Status      int32        `json:"status"`
	Buffer      bytes.Buffer `json:"buffer,omitempty"`
	Msg         string       `json:"msg,omitempty"`
	Txid        string       `json:"txid"`
	Function    string       `json:"function,omitempty"`
	Mspid       string       `json:"mspid,omitempty"`
	CertIssuer  string       `json:"certIssuer,omitempty"`
	CertSubject string       `json:"certSubject,omitempty"`
	Assets      int          `json:"assets,omitempty"`
	WrongItem   int          `json:"wrongItem,omitempty"`
}

func systemErrorResponse(msg string, wrongItem ...int) Response {
	return errorResponse(msg, INTERNAL_SERVER_ERROR, wrongItem)
}

func clientErrorResponse(msg string, wrongItem ...int) Response {
	return errorResponse(msg, BAD_REQUEST, wrongItem)
}

func forbiddenErrorResponse(msg string, wrongItem ...int) Response {
	return errorResponse(msg, FORBIDDEN, wrongItem)
}

func errorResponse(msg string, status int32, wrongItem []int) Response {
	var response Response
	response.Status = status
	response.Msg = msg
	if len(wrongItem) > 0 {
		response.WrongItem = wrongItem[0]
	}
	return response
}

func successResponse(msg string, assets int) Response {
	var response Response
	response.Status = OK
	response.Assets = assets
	response.Msg = msg
	return response
}

func successResponseWithBuffer(buffer *bytes.Buffer) Response {
	var response Response
	response.Status = OK
	response.Buffer = *buffer
	return response
}

func (response *Response) peerResponse(ctx Ctx) peer.Response {
	if response.isOk() && response.Buffer.Len() > 0 {
		return shim.Success(response.Buffer.Bytes())
	} else {
		response.Buffer.Reset()
		response.Txid = ctx.txid
		if ctx.verboseMode || response.isError() {
			response.Function = ctx.function
			response.Mspid = ctx.mspid
			response.CertIssuer = ctx.certIssuer
			response.CertSubject = ctx.certSubject
		} else {
			response.Msg = ""
			response.Function = ""
			response.Mspid = ""
			response.CertIssuer = ""
			response.CertSubject = ""
		}
		responseAsBytes, _ := json.Marshal(response)
		if response.isError() {
			return shim.Error(string(responseAsBytes))
		} else {
			return shim.Success(responseAsBytes)
		}
	}
}

func (r *Response) isError() bool {
	return !r.isOk()
}

func (r *Response) isOk() bool {
	return r.Status == OK || r.Status == 0
}
