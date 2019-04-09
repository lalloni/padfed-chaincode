package fabric

import (
	"bytes"
)

type responseCode int32

const (
	Unknown     responseCode = 0
	OK          responseCode = 200
	BadRequest  responseCode = 400
	Forbidden   responseCode = 403
	NotFound    responseCode = 404
	ServerError responseCode = 500
)

type Response struct {
	Status      responseCode  `json:"status"`
	Buffer      *bytes.Buffer `json:"-"`
	Msg         string        `json:"msg,omitempty"`
	Txid        string        `json:"txid"`
	Function    string        `json:"function,omitempty"`
	Mspid       string        `json:"mspid,omitempty"`
	CertIssuer  string        `json:"certIssuer,omitempty"`
	CertSubject string        `json:"certSubject,omitempty"`
	Assets      int           `json:"assets,omitempty"`
	WrongItem   int           `json:"wrongItem,omitempty"`
}

func (r *Response) IsOK() bool {
	return r.Status == OK || r.Status == Unknown
}

func NotFoundErrorResponse() *Response {
	return ErrorResponse("", NotFound)
}

func SystemErrorResponse(msg string, wrongItem ...int) *Response {
	return ErrorResponse(msg, ServerError, wrongItem...)
}

func ClientErrorResponse(msg string, wrongItem ...int) *Response {
	return ErrorResponse(msg, BadRequest, wrongItem...)
}

func ForbiddenErrorResponse(msg string, wrongItem ...int) *Response {
	return ErrorResponse(msg, Forbidden, wrongItem...)
}

func ErrorResponse(msg string, status responseCode, wrongItem ...int) *Response {
	var response Response
	response.Status = status
	response.Msg = msg
	if len(wrongItem) > 0 {
		response.WrongItem = wrongItem[0]
	}
	return &response
}

func SuccessResponse(msg string, assets int) *Response {
	var response Response
	response.Status = OK
	response.Assets = assets
	response.Msg = msg
	return &response
}

func SuccessResponseWithBuffer(buffer *bytes.Buffer) *Response {
	var response Response
	response.Status = OK
	response.Buffer = buffer
	return &response
}
