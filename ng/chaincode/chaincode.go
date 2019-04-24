package chaincode

import (
	"bytes"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func New(log *shim.ChaincodeLogger, r router.Router) shim.Chaincode {
	return &cc{
		log:    log,
		router: r,
	}
}

type cc struct {
	log    *shim.ChaincodeLogger
	router router.Router
}

func (c *cc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	ctx := context.New(stub)
	handle := c.router.InitHandler()
	if handle != nil {
		return c.response(ctx, handle(ctx))
	}
	return c.response(ctx, response.OK(nil))
}

func (c *cc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	ctx := context.New(stub)
	handle := c.router.InvokeHandler(ctx.Function())
	if handle == nil {
		handle = handler.NotImplementedHandler
	}
	return c.response(ctx, handle(ctx))
}

func (c *cc) response(ctx *context.Context, r *response.Response) peer.Response {
	var payload []byte
	if r.Status >= shim.ERRORTHRESHOLD {
		if r.Payload == nil {
			r.Payload = &response.Payload{}
		}
		r.Payload.Transaction = &response.Transaction{ID: ctx.Stub.GetTxID(), Function: ctx.Function()}
		mspid, err := ctx.ClientMSPID()
		if err != nil {
			c.log.Errorf("getting MSPID: %v", err)
		}
		var subject, issuer string
		cert, err := ctx.ClientCertificate()
		if err != nil {
			c.log.Errorf("getting client certificate: %v", err)
		} else {
			subject = cert.Subject.String()
			issuer = cert.Issuer.String()
		}
		r.Payload.Client = &response.Client{MSPID: mspid, Subject: subject, Issuer: issuer}
	}
	if r.Payload != nil {
		b := &bytes.Buffer{}
		err := json.NewEncoder(b).Encode(r.Payload)
		if err != nil {
			return c.response(ctx, response.Error("encoding response payload: %v", err))
		}
		payload = b.Bytes()
	}
	return peer.Response{
		Status:  r.Status,
		Message: r.Message,
		Payload: payload,
	}
}
