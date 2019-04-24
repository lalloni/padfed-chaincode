package context

import (
	"crypto/x509"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func New(stub shim.ChaincodeStubInterface) *Context {
	return &Context{
		Stub:  stub,
		Store: store.New(stub),
	}
}

type Context struct {
	Stub        shim.ChaincodeStubInterface
	Store       store.Store
	clientid    cid.ClientIdentity
	clientcrt   *x509.Certificate
	clientmspid string
	function    Function
}

func (ctx *Context) ClientIdentity() (cid.ClientIdentity, error) {
	if ctx.clientid == nil {
		id, err := cid.New(ctx.Stub)
		if err != nil {
			return nil, errors.Wrap(err, "creating client identity")
		}
		ctx.clientid = id
	}
	return ctx.clientid, nil
}

func (ctx *Context) ClientCertificate() (*x509.Certificate, error) {
	if ctx.clientcrt == nil {
		id, err := ctx.ClientIdentity()
		if err != nil {
			return nil, err
		}
		cert, err := id.GetX509Certificate()
		if err != nil {
			return nil, errors.Wrap(err, "getting client certificate")
		}
		ctx.clientcrt = cert
	}
	return ctx.clientcrt, nil
}

func (ctx *Context) ClientMSPID() (string, error) {
	if ctx.clientmspid == "" {
		id, err := ctx.ClientIdentity()
		if err != nil {
			return "", err
		}
		mspid, err := id.GetMSPID()
		if err != nil {
			return "", errors.Wrap(err, "getting client mspid")
		}
		ctx.clientmspid = mspid
	}
	return ctx.clientmspid, nil
}

func (ctx *Context) Function() Function {
	if ctx.function == Function("") {
		ctx.function = Function(string(ctx.Stub.GetArgs()[0]))
	}
	return ctx.function
}

func (ctx *Context) ArgBytes(n int) ([]byte, error) {
	args := ctx.Stub.GetArgs()
	if len(args) < n+1 {
		return nil, errors.Errorf("argument %d required", n)
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
