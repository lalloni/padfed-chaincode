package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

type formatType int

const (
	JSON     formatType = 1
	PROTOBUF formatType = 2
)

// Persona asset
type Persona struct {
	CUIT          uint64      `protobuf:"varint,1,name=cuit,proto3" json:"cuit"`
	Nombre        string      `protobuf:"bytes,2,opt,name=nombre,proto3" json:"nombre,omitempty"`
	Apellido      string      `protobuf:"bytes,3,opt,name=apellido,proto3" json:"apellido,omitempty"`
	RazonSocial   string      `protobuf:"bytes,4,opt,name=razon_social,proto3" json:"razonSocial,omitempty"`
	Tipo          string      `protobuf:"bytes,5,name=tipo,proto3" json:"tipo,omitempty"`
	Estado        string      `protobuf:"bytes,6,name=estado,proto3" json:"estado,omitempty"`
	FormaJuridica int32       `protobuf:"varint,7,opt,name=forma_juridica,proto3" json:"formaJuridica,omitempty"`
	TipoDoc       int32       `protobuf:"varint,8,opt,name=tipo_doc,proto3" json:"tipoDoc,omitempty"`
	Doc           string      `protobuf:"bytes,9,opt,name=doc,proto3" json:"doc,omitempty"`
	Sexo          string      `protobuf:"bytes,10,opt,name=sexo,proto3" json:"sexo,omitempty"`
	MesCierre     int32       `protobuf:"varint,11,opt,name=mes_cierre,proto3" json:"mesCierre,omitempty"`
	Nacimiento    string      `protobuf:"bytes,12,opt,name=nacimiento,proto3" json:"nacimiento,omitempty"`
	Fallecimiento string      `protobuf:"bytes,13,opt,name=fallecimiento,proto3" json:"fallecimiento,omitempty"`
	Inscripcion   string      `protobuf:"bytes,14,opt,name=inscripcion,proto3" json:"inscripcion,omitempty"`
	FechaCierre   string      `protobuf:"bytes,15,opt,name=fecha_cierre,proto3" json:"fechaCierre,omitempty"`
	NuevaCUIT     uint64      `protobuf:"varint,16,opt,name=nueva_cuit,proto3" json:"nuevaCuit,omitempty"`
	Materno       string      `protobuf:"bytes,17,opt,name=materno,proto3" json:"materno,omitempty"`
	Pais          string      `protobuf:"bytes,18,opt,name=pais,proto3" json:"pais,omitempty"`
	CH            string      `protobuf:"bytes,19,opt,name=ch,proto3" json:"ch,omitempty"`
	DS            string      `protobuf:"bytes,20,opt,name=ds,proto3" json:"ds,omitempty"`
	Impuestos     []*Impuesto `protobuf:"group,21,rep,name=impuestos,proto3" json:"impuestos,omitempty"`
}

func (m *Persona) Reset()         { *m = Persona{} }
func (m *Persona) String() string { return proto.CompactTextString(m) }
func (*Persona) ProtoMessage()    {}
func (m *Persona) GetImpuestos() []*Impuesto {
	if m != nil {
		return m.Impuestos
	}
	return nil
}

// Impuesto asset
type Impuesto struct {
	Impuesto    int32  `protobuf:"varint,1,name=impuesto,proto3" json:"impuesto"`
	IDOrganismo int32  `protobuf:"varint,2,opt,name=id_org,proto3" json:"idOrg,omitempty"`
	Inscripcion string `protobuf:"bytes,3,opt,name=inscripcion,proto3" json:"inscripcion,omitempty"`
	Periodo     int32  `protobuf:"varint,4,opt,name=periodo,proto3" json:"periodo"`
	Estado      string `protobuf:"bytes,5,opt,name=estado,proto3" json:"estado"`
	IDTxc       uint64 `protobuf:"varint,6,opt,name=id_txc,proto3" json:"idTxc,omitempty"`
	DS          string `protobuf:"bytes,7,opt,name=ds,proto3" json:"ds,omitempty"`
	Motivo      string `protobuf:"bytes,8,opt,name=motivo,proto3" json:"motivo,omitempty"`
	Dia         int32  `protobuf:"varint,9,opt,name=dia,proto3" json:"dia,omitempty"`
}

func (m *Impuesto) Reset()         { *m = Impuesto{} }
func (m *Impuesto) String() string { return proto.CompactTextString(m) }
func (*Impuesto) ProtoMessage()    {}

type Personas struct {
	Personas []*Persona `protobuf:"group,1,rep,name=personas,proto3" json:"personas"`
}

func (m *Personas) Reset()         { *m = Personas{} }
func (m *Personas) String() string { return proto.CompactTextString(m) }
func (*Personas) ProtoMessage()    {}
func (m *Personas) GetPersonas() []*Persona {
	if m != nil {
		return m.Personas
	}
	return nil
}

// Impuestos es una estructura solo usada para poder extraer un array de impuestos.
type Impuestos struct {
	Impuestos []*Impuesto `json:"impuestos"`
}

// ParamImpuesto asset
type ParamImpuesto struct {
	Impuesto    int32  `json:"impuesto"`
	IDOrganismo int32  `json:"idOrg"`
	TipoRegimen string `json:"tipoRegimen"`
	Nombre      string `json:"nombre"`
	FechaDesde  string `json:"fechaDesde"`
	FechaHasta  string `json:"fechaHasta"`
}

// TXConfirmable asset
type TXConfirmable struct {
	IDTxc              uint64    `json:"idTxc"`
	TipoTxc            int       `json:"tipoTxc"`
	NombreTxc          int       `json:"nombreTxc"`
	IDOrganismo        int       `json:"idOrg"`
	CUIT               uint64    `json:"cuit"`
	AssetType          string    `json:"assetType"`
	AssetValue         string    `json:"assetValue"`
	FechaHoraTxc       time.Time `json:"fechahoraTxc"`
	FechaHoraRespuesta time.Time `json:"fechahoraRespuesta"`
	TipoRespuesta      int       `json:"tipoRespuesta"`
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
var FIND_VERBOSE_REGEXP = *regexp.MustCompile(`^(.*)(\?v)$`)

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
		case "putPersona",
			"putPersonaProto",
			"putPersonas":
			if err := checkClientID(ctx); err.isError() {
				return err.peerResponse(ctx)
			}
		}
	}
	switch ctx.function {
	case "putPersona":
		r = s.putPersona(APIstub, ctx.args, JSON)
	case "putPersonaProto":
		r = s.putPersona(APIstub, ctx.args, PROTOBUF)
	case "putPersonas":
		r = s.putPersonas(APIstub, ctx.args, JSON)
	case "putPersonasProto":
		r = s.putPersonas(APIstub, ctx.args, PROTOBUF)
	case "putParamImpuestos":
		r = s.putParamImpuestos(APIstub, ctx.args)
	case "createTxConfirmable":
		r = s.createTxConfirmable(APIstub, ctx.args)
	case "responseTxConfirmable":
		r = s.responseTxConfirmable(APIstub, ctx.args)
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
	case "queryTxConfirmables":
		r = s.queryTxConfirmables(APIstub, ctx.args)
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
		ctx.verboseMode = true
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
