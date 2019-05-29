package chaincode

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/impuestos"
)

type Ctx struct {
	version     string
	verboseMode bool
	// current data transaction
	txid     string
	function string
	args     []string
	// Datos de clientIdentity
	mspid       string
	certIssuer  string
	certSubject string
}

func New(log *shim.ChaincodeLogger, version string, testing bool) shim.Chaincode {
	return &padfedcc{
		handlers: BuildHandlers(version, testing),
		log:      log,
		testing:  testing,
		version:  version,
	}
}

type padfedcc struct {
	testing  bool
	log      *shim.ChaincodeLogger
	version  string
	handlers Handlers
}

func (s *padfedcc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	ctx, r := setContext(stub, s.version, s.testing)
	if !r.IsOK() {
		return peerResponse(ctx, r)
	}
	if !s.testing {
		if ctx.mspid != AFIP {
			return peerResponse(ctx, fabric.ForbiddenErrorResponse("MSPID must be AFIP"))
		}
	}
	r = impuestos.LoadInitial(stub)
	return peerResponse(ctx, r)
}

func (s *padfedcc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	var ctx *Ctx
	var r *fabric.Response

	if ctx, r = setContext(stub, s.version, s.testing); !r.IsOK() {
		return peerResponse(ctx, r)
	}

	s.log.Infof("processing invocation of %s(%v) on transaction %s", ctx.function, ctx.args, ctx.txid)

	if handler, ok := s.handlers[ctx.function]; ok {
		r = handler(stub, ctx.args)
	} else {
		r = fabric.ClientErrorResponse("unknown chaincode function name " + ctx.function)
	}

	return peerResponse(ctx, r)
}

func setContext(stub shim.ChaincodeStubInterface, version string, testing bool) (*Ctx, *fabric.Response) {
	ctx := &Ctx{}
	ctx.txid = stub.GetTxID()
	ctx.version = version
	ctx.function, ctx.args = stub.GetFunctionAndParameters()
	ff := strings.SplitN(ctx.function, "?", 2)
	if len(ff) > 1 {
		ctx.function = ff[0]
		qry, err := url.ParseQuery(ff[1])
		if err != nil {
			return &Ctx{}, fabric.SystemErrorResponse(fmt.Sprintf("Error parsing function options: %v", err))
		}
		if qry.Get("verbose") == "true" {
			ctx.verboseMode = true
		}
	}
	if !testing {
		// Get the client ID object
		clientIdentity, err := cid.New(stub)
		if err != nil {
			return &Ctx{}, fabric.SystemErrorResponse("Error at Get the client ID object [cid.New(stub)]")
		}
		mspid, err := clientIdentity.GetMSPID()
		if err != nil {
			return &Ctx{}, fabric.SystemErrorResponse("Error at Get the client ID object [GetMSPID()]")
		}
		ctx.mspid = mspid

		x509Certificate, err := clientIdentity.GetX509Certificate()
		if err != nil {
			return &Ctx{}, fabric.SystemErrorResponse("Error at Get the x509Certificate object [GetX509Certificate()]")
		}
		ctx.certSubject = x509Certificate.Subject.String()
		ctx.certIssuer = x509Certificate.Issuer.String()
	}
	return ctx, &fabric.Response{}
}

func peerResponse(ctx *Ctx, response *fabric.Response) peer.Response {
	if response.IsOK() && response.Buffer != nil {
		return shim.Success(response.Buffer.Bytes())
	}
	response.Txid = ctx.txid
	if ctx.verboseMode || !response.IsOK() {
		response.Version = ctx.version
		response.Function = ctx.function
		response.Mspid = ctx.mspid
		response.CertIssuer = ctx.certIssuer
		response.CertSubject = ctx.certSubject
	} else {
		response.Msg = ""
		response.Function = ""
		response.Version = ""
		response.Mspid = ""
		response.CertIssuer = ""
		response.CertSubject = ""
	}
	responseAsBytes, _ := json.Marshal(response)
	if response.IsOK() {
		return shim.Success(responseAsBytes)
	}
	return shim.Error(string(responseAsBytes)) // TODO (pil): esto devuelve un json tuneleado en una string en un json... wtf?
}
