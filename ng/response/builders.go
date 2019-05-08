package response

import (
	"fmt"
	"net/http"
)

func OK(result interface{}) *Response {
	return StatusWithResult(http.StatusOK, result)
}

func NotFound() *Response {
	return Status(http.StatusNotFound)
}

func Forbidden(msg string, args ...interface{}) *Response {
	return StatusWithMessage(http.StatusForbidden, msg, args...)
}

func BadRequest(msg string, args ...interface{}) *Response {
	return StatusWithMessage(http.StatusBadRequest, msg, args...)
}

func BadRequestWithFault(fault interface{}) *Response {
	return StatusWithFault(http.StatusBadRequest, fault)
}

func Error(msg string, args ...interface{}) *Response {
	return StatusWithMessage(http.StatusInternalServerError, msg, args...)
}

func NotImplemented() *Response {
	return Status(http.StatusNotImplemented)
}

func StatusWithResult(status int32, result interface{}) *Response {
	return &Response{
		Status: status,
		Payload: &Payload{
			Content: result,
		},
	}
}

func StatusWithFault(status int32, fault interface{}) *Response {
	return &Response{
		Status: status,
		Payload: &Payload{
			Fault: fault,
		},
	}
}

func StatusWithMessage(status int32, msg string, args ...interface{}) *Response {
	return &Response{
		Status:  status,
		Message: fmt.Sprintf(msg, args...),
	}
}

func Status(status int32) *Response {
	return &Response{
		Status: status,
	}
}
