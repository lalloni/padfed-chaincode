package deprecated

import (
	"fmt"

	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func Adapter(old chaincode.Handler, use string) handler.Handler {
	w := fmt.Sprintf("This function is DEPRECATED and will be REMOVED! Please migrate to using function %q.", use)
	return func(ctx *context.Context) *response.Response {
		ctx.Logger().Notice("deprecated function called instead of %q", use)
		c, res := chaincode.SetContext(ctx.Stub, ctx.Version(), false)
		if !res.IsOK() {
			return response.Direct(warn2(chaincode.PeerResponse(c, warn1(res, w)), w))
		}
		_, args := ctx.Stub.GetFunctionAndParameters()
		res = old(ctx.Stub, args)
		return response.Direct(warn2(chaincode.PeerResponse(c, warn1(res, w)), w))
	}
}

func warn1(p *fabric.Response, msg string) *fabric.Response {
	p.Warning = msg
	return p
}

func warn2(p peer.Response, msg string) peer.Response {
	if p.Message == "" {
		p.Message = msg
	}
	return p
}
