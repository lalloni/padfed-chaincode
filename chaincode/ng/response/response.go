package response

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric/protos/peer"
)

func Status(status int32, msg string, args ...interface{}) peer.Response {
	return peer.Response{
		Status:  status,
		Message: fmt.Sprintf(msg, args...),
	}
}

func BadRequest(msg string, args ...interface{}) peer.Response {
	return Status(http.StatusBadRequest, msg, args...)
}

func Forbidden(msg string, args ...interface{}) peer.Response {
	return Status(http.StatusForbidden, msg, args...)
}

func NotFound(msg string, args ...interface{}) peer.Response {
	return Status(http.StatusNotFound, msg, args...)
}

func Error(msg string, args ...interface{}) peer.Response {
	return Status(http.StatusInternalServerError, msg, args...)
}

func ErrorWrap(err error, msg string, args ...interface{}) peer.Response {
	r := Status(http.StatusInternalServerError, msg, args...)
	r.Message += ": " + err.Error()
	return r
}

func Success(bs []byte) peer.Response {
	return peer.Response{
		Status:  http.StatusOK,
		Payload: bs,
	}
}
