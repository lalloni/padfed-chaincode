package ng

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/acl"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/support"
)

func New(log *shim.ChaincodeLogger) shim.Chaincode {
	return &cc{
		Log:      log,
		Handlers: Handlers(),
	}
}

type cc struct {
	Log      *shim.ChaincodeLogger
	Handlers support.HandlerMap
}

func (c *cc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	ctx, err := support.NewContext(stub)
	if err != nil {
		return response.ErrorWrap(err, "creating new context")
	}
	if ctx.MSPID != acl.AFIPMSPID {
		return response.Forbidden("MSPID is not %q", acl.AFIPMSPID)
	}
	return response.Success(nil)
}

func (c *cc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	ctx, err := support.NewContext(stub)
	if err != nil {
		return response.ErrorWrap(err, "creating new context")
	}
	function := ctx.Function()
	return c.Handlers[function](ctx)
}
