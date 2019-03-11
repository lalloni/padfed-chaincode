package chaincode

import (
	"encoding/json"
	"regexp"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/impuestos"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Ctx struct {
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

func New(log *shim.ChaincodeLogger, testing bool) shim.Chaincode {
	return &padfedcc{
		handlers: BuildHandlers(),
		log:      log,
		testing:  testing,
	}
}

type padfedcc struct {
	testing  bool
	log      *shim.ChaincodeLogger
	handlers Handlers
}

var verboseRegexp = regexp.MustCompile(`^(.*)(\?verbose=)(true|false)$`)

func (s *padfedcc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	ctx, r := setContext(stub, s.testing)
	if !r.IsOK() {
		return peerResponse(ctx, r)
	}
	if !s.testing {
		r = checkClientID(ctx)
		if !r.IsOK() {
			return peerResponse(ctx, r)
		}
	}
	r = impuestos.LoadInitial(stub)
	return peerResponse(ctx, r)
}

func (s *padfedcc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	var ctx *Ctx
	var r *fabric.Response

	if ctx, r = setContext(stub, s.testing); !r.IsOK() {
		return peerResponse(ctx, r)
	}

	s.log.Infof("processing invocation of %s(%v) on transaction %s", ctx.function, ctx.args, ctx.txid)

	if !s.testing {
		switch ctx.function {
		case "putPersona", "putPersonas":
			if err := checkClientID(ctx); !err.IsOK() {
				return peerResponse(ctx, err)
			}
		}
	}

	if handler, ok := s.handlers[ctx.function]; ok {
		r = handler(stub, ctx.args)
	} else {
		r = fabric.ClientErrorResponse("Invalid Smart Contract function name " + ctx.function)
	}

	return peerResponse(ctx, r)
}

func setContext(stub shim.ChaincodeStubInterface, testing bool) (*Ctx, *fabric.Response) {
	ctx := &Ctx{}
	ctx.txid = stub.GetTxID()
	ctx.function, ctx.args = stub.GetFunctionAndParameters()
	// Check for verbose mode
	res := verboseRegexp.FindStringSubmatch(ctx.function)
	if len(res) != 0 {
		ctx.function = res[1]
		if res[3] == "true" {
			ctx.verboseMode = true
		} else {
			ctx.verboseMode = false
		}
	} else {
		ctx.verboseMode = false
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

func checkClientID(ctx *Ctx) *fabric.Response {
	if ctx.mspid != "AFIP" {
		return fabric.ForbiddenErrorResponse("mspid [" + ctx.mspid + "] - La funcion [" + ctx.function + "] solo puede ser invocada por AFIP")
	}
	return &fabric.Response{}
}

func peerResponse(ctx *Ctx, response *fabric.Response) peer.Response {
	if response.IsOK() && response.Buffer.Len() > 0 {
		return shim.Success(response.Buffer.Bytes())
	}
	response.Txid = ctx.txid
	if ctx.verboseMode || !response.IsOK() {
		response.Function = ctx.function
		response.Mspid = ctx.mspid
		response.CertIssuer = ctx.certIssuer
		response.CertSubject = ctx.certSubject
	} else {
		response.Msg = ""
		response.Function = ""
		response.Mspid = ""
		response.CertIssuer = ""
		response.CertSubject = ""
	}
	responseAsBytes, _ := json.Marshal(response)
	if response.IsOK() {
		return shim.Success(responseAsBytes)
	}
	return shim.Error(string(responseAsBytes))

}
