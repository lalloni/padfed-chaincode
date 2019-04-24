package support

import (
	"crypto/x509"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func NewContext(stub shim.ChaincodeStubInterface) (*Context, error) {
	id, err := cid.New(stub)
	if err != nil {
		return nil, errors.Wrap(err, "creating client identity")
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		return nil, errors.Wrap(err, "getting client mspid")
	}
	cert, err := id.GetX509Certificate()
	if err != nil {
		return nil, errors.Wrap(err, "getting client certificate")
	}
	return &Context{
		Stub:        stub,
		Store:       store.New(stub),
		Identity:    id,
		MSPID:       mspid,
		Certificate: cert,
	}, nil
}

type Context struct {
	Stub        shim.ChaincodeStubInterface
	Store       store.Store
	Identity    cid.ClientIdentity
	Certificate *x509.Certificate
	MSPID       string
}

func (ctx *Context) Function() string {
	return string(ctx.Stub.GetArgs()[0])
}

func (ctx *Context) ArgBytes(n int) ([]byte, error) {
	args := ctx.Stub.GetArgs()
	if len(args) < n+1 {
		return nil, errors.New("insufficient arguments")
	}
	return args[n], nil
}

func (ctx *Context) ArgString(n int) (string, error) {
	bs, err := ctx.ArgBytes(n)
	return string(bs), err
}

func (ctx *Context) ArgInt(n int) (int, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(string(bs), 10, strconv.IntSize)
	return int(v), err
}

func (ctx *Context) ArgUint(n int) (uint, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(string(bs), 10, strconv.IntSize)
	return uint(v), err
}
