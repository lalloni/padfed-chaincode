package context

import (
	"crypto/x509"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/logging"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func New(stub shim.ChaincodeStubInterface, path ...string) *Context {
	return &Context{
		path:  path,
		Stub:  stub,
		Store: store.New(stub),
	}
}

type Context struct {
	Stub        shim.ChaincodeStubInterface
	Store       store.Store
	path        []string
	clientid    cid.ClientIdentity
	clientcrt   *x509.Certificate
	clientmspid string
	function    string
}

func (ctx *Context) Logger(path ...string) *shim.ChaincodeLogger {
	return logging.ChaincodeLogger(append(ctx.path, path...)...)
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

func (ctx *Context) Function() string {
	if ctx.function == "" {
		ctx.function = string(ctx.Stub.GetArgs()[0])
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

func (ctx *Context) ArgInt64(n int) (int64, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(bs), 10, 64)
}

func (ctx *Context) ArgUint64(n int) (uint64, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(bs), 10, 64)
}
