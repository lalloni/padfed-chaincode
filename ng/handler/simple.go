package handler

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func SuccessHandler(ctx *context.Context) *response.Response {
	return response.OK(nil)
}

func EchoHandler(ctx *context.Context) *response.Response {
	data, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest(err.Error())
	}
	return response.OK(data)
}

func VersionHandler(version string) Handler {
	return func(ctx *context.Context) *response.Response {
		return response.OK(version)
	}
}

func NotImplementedHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}
