package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Persona asset
type Persona struct {
	CUIT          uint64      `json:"cuit"`
	Nombre        string      `json:"nombre,omitempty"`
	Apellido      string      `json:"apellido,omitempty"`
	RazonSocial   string      `json:"razonSocial,omitempty"`
	Tipo          string      `json:"tipo,omitempty"`
	Estado        string      `json:"estado,omitempty"`
	FormaJuridica int32       `json:"formaJuridica,omitempty"`
	TipoDoc       int32       `json:"tipoDoc,omitempty"`
	Doc           string      `json:"doc,omitempty"`
	Sexo          string      `json:"sexo,omitempty"`
	MesCierre     int32       `json:"mesCierre,omitempty"`
	Nacimiento    string      `json:"nacimiento,omitempty"`
	Fallecimiento string      `json:"fallecimiento,omitempty"`
	Inscripcion   string      `json:"inscripcion,omitempty"`
	FechaCierre   string      `json:"fechaCierre,omitempty"`
	NuevaCUIT     uint64      `json:"nuevaCuit,omitempty"`
	Materno       string      `json:"materno,omitempty"`
	Pais          string      `json:"pais,omitempty"`
	CH            []string    `json:"ch,omitempty"`
	DS            string      `json:"ds,omitempty"`
	Impuestos     []*Impuesto `json:"impuestos,omitempty"`
}

// Impuesto asset
type Impuesto struct {
	Impuesto    int32  `json:"impuesto"`
	IDOrganismo int32  `json:"idOrg,omitempty"`
	Inscripcion string `json:"inscripcion,omitempty"`
	Periodo     int32  `json:"periodo"`
	Estado      string `json:"estado"`
	IDTxc       uint64 `json:"idTxc,omitempty"`
	DS          string `json:"ds,omitempty"`
	Motivo      string `json:"motivo,omitempty"`
	Dia         int32  `json:"dia,omitempty"`
}

type Personas struct {
	Personas []*Persona `json:"personas"`
}

// Impuestos es una estructura solo usada para poder extraer un array de impuestos.
type Impuestos struct {
	Impuestos []*Impuesto `json:"impuestos"`
}

// ParamImpuesto asset
type ParamImpuesto struct {
	Impuesto    int32  `json:"impuesto"`
	IDOrg       int32  `json:"idOrg"`
	TipoRegimen string `json:"tipoRegimen"`
	Nombre      string `json:"nombre"`
	FechaDesde  string `json:"fechaDesde"`
	FechaHasta  string `json:"fechaHasta"`
}

// Context
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

// SmartContract Agrupador de funciones
type SmartContract struct {
	isModeTest bool
}

var logger = shim.NewLogger("rut-afipcc")
var FIND_VERBOSE_REGEXP = *regexp.MustCompile(`^(.*)(\?verbose=)(true|false)$`)

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	var ctx Ctx
	var r Response
	if ctx, r = setContext(APIstub, s.isModeTest); r.isError() {
		return r.peerResponse(ctx)
	}
	if !s.isModeTest {
		if r = checkClientID(ctx); r.isError() {
			return r.peerResponse(ctx)
		}
	}
	r = s.setInitImpuestos(APIstub)
	return r.peerResponse(ctx)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	var ctx Ctx
	var r Response
	if ctx, r = setContext(APIstub, s.isModeTest); r.isError() {
		return r.peerResponse(ctx)
	}
	log.Print("=================================================================")
	log.Print("TxID [" + ctx.txid + "]")
	log.Print("Function [" + ctx.function + "] args [" + strings.Join(ctx.args, " ") + "]")

	if !s.isModeTest {
		switch ctx.function {
		case "putPersona", "putPersonas":
			if err := checkClientID(ctx); err.isError() {
				return err.peerResponse(ctx)
			}
		}
	}

	switch ctx.function {
	case "putPersona":
		r = s.putPersona(APIstub, ctx.args)
	case "putPersonas":
		r = s.putPersonas(APIstub, ctx.args)
	case "putParamImpuestos":
		r = s.putParamImpuestos(APIstub, ctx.args)
	case "delPersona":
		r = s.delPersona(APIstub, ctx.args)
	case "delPersonasByRange":
		r = s.delPersonasByRange(APIstub, ctx.args)
	case "deleteAll":
		r = s.deleteByKeyRange(APIstub, []string{"", ""})
	case "deleteByKeyRange":
		r = s.deleteByKeyRange(APIstub, ctx.args)
	case "delParamImpuestosAll":
		r = s.delParamImpuestosAll(APIstub)
	case "delPersonaAssets":
		r = s.delPersonaAssets(APIstub, ctx.args)
	case "putPersonaImpuestos":
		r = s.putPersonaImpuestos(APIstub, ctx.args)
	case "queryPersona":
		r = s.queryPersona(APIstub, ctx.args)
	case "queryPersonasByRange":
		r = s.queryPersonasByRange(APIstub, ctx.args)
	case "queryPersonaImpuestos":
		r = s.queryPersonaImpuestos(APIstub, ctx.args)
	case "queryAllPersona":
		r = s.queryAllPersona(APIstub)
	case "queryByKey":
		r = s.queryByKey(APIstub, ctx.args)
	case "queryByKeyRange":
		r = s.queryByKeyRange(APIstub, ctx.args)
	case "queryParamImpuestos":
		r = s.queryParamImpuestos(APIstub)
	case "queryHistory":
		r = s.queryHistory(APIstub, ctx.args)
	default:
		r = clientErrorResponse("Invalid Smart Contract function name " + ctx.function)
	}
	return r.peerResponse(ctx)
}

func setContext(APIstub shim.ChaincodeStubInterface, isModeTest bool) (Ctx, Response) {
	var ctx Ctx
	ctx.txid = APIstub.GetTxID()
	ctx.function, ctx.args = APIstub.GetFunctionAndParameters()
	// Check for verbose mode
	res := FIND_VERBOSE_REGEXP.FindStringSubmatch(ctx.function)
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
	if !isModeTest {
		// Get the client ID object
		clientIdentity, err := cid.New(APIstub)
		if err != nil {
			return Ctx{}, systemErrorResponse("Error at Get the client ID object [cid.New(APIstub)]")
		}
		mspid, err := clientIdentity.GetMSPID()
		if err != nil {
			return Ctx{}, systemErrorResponse("Error at Get the client ID object [GetMSPID()]")
		}
		ctx.mspid = mspid

		x509Certificate, err := clientIdentity.GetX509Certificate()
		if err != nil {
			return Ctx{}, systemErrorResponse("Error at Get the x509Certificate object [GetX509Certificate()]")
		}
		ctx.certSubject = x509Certificate.Subject.String()
		ctx.certIssuer = x509Certificate.Issuer.String()
	}
	return ctx, Response{}
}

func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("Ledger inicializado")
	return shim.Success(nil)
}

func checkClientID(ctx Ctx) Response {
	if ctx.mspid != "AFIP" {
		return forbiddenErrorResponse("mspid [" + ctx.mspid + "] - La funcion [" + ctx.function + "] solo puede ser invocada por AFIP")
	}
	return Response{}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	// LogDebug, LogInfo, LogNotice, LogWarning, LogError, LogCritical (Default: LogDebug)
	logger.SetLevel(shim.LogDebug)

	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)
	if err := shim.Start(new(SmartContract)); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
