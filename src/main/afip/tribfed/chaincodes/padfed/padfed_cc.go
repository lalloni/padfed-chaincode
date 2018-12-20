package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
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
	CUIT               uint64      `protobuf:"varint,1,name=cuit,proto3" json:"cuit"`
	Nombre             string      `protobuf:"bytes,2,opt,name=nombre,proto3" json:"nombre,omitempty"`
	Apellido           string      `protobuf:"bytes,3,opt,name=apellido,proto3" json:"apellido,omitempty"`
	RazonSocial        string      `protobuf:"bytes,4,opt,name=razon_social,proto3" json:"razonSocial,omitempty"`
	TipoPersona        string      `protobuf:"bytes,5,name=tipo_persona,proto3" json:"tipoPersona"`
	EstadoCUIT         string      `protobuf:"bytes,6,name=estado_cuit,proto3" json:"estadoCuit"`
	IDFormaJuridica    int32       `protobuf:"varint,7,opt,name=tipo_persona,proto3" json:"idFormaJuridica,omitempty"`
	TipoDoc            int32       `protobuf:"varint,8,opt,name=tipo_doc,proto3" json:"tipoDoc,omitempty"`
	Documento          string      `protobuf:"bytes,9,opt,name=documento,proto3" son:"documento,omitempty"`
	Sexo               string      `protobuf:"bytes,10,opt,name=sexo,proto3" son:"sexo,omitempty"`
	MesCierre          int32       `protobuf:"varint,11,opt,name=mes_cierre,proto3" son:"mesCierre,omitempty"`
	FechaNacimiento    string      `protobuf:"bytes,12,opt,name=fecha_nacimiento,proto3" son:"fechaNacimiento,omitempty"`
	FechaFallecimiento string      `protobuf:"bytes,13,opt,name=fecha_fallecimiento,proto3" son:"fechaFallecimiento,omitempty"`
	FechaInscripcion   string      `protobuf:"bytes,14,opt,name=fecha_inscripcion,proto3" son:"fechaInscripcion,omitempty"`
	FechaCierre        string      `protobuf:"bytes,15,opt,name=fecha_cierre,proto3" son:"fechaCierre,omitempty"`
	NuevaCUIT          uint64      `protobuf:"varint,16,opt,name=nueva_cuit,proto3" son:"nuevaCuit,omitempty"`
	Impuestos          []*Impuesto `protobuf:"group,17,rep,name=impuestos,proto3" json:"impuestos,omitempty"`
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
	IDImpuesto       int32  `protobuf:"varint,1,name=id_impuesto,proto3" json:"idImpuesto"`
	IDOrganismo      int32  `protobuf:"varint,2,opt,name=id_org,proto3" json:"idOrg,omitempty"`
	FechaInscripcion string `protobuf:"bytes,3,opt,name=fecha_inscripcion,proto3" json:"fechaInscripcion,omitempty"`
	Periodo          int32  `protobuf:"varint,4,opt,name=periodo,proto3" json:"periodo,omitempty"`
	Estado           string `protobuf:"bytes,5,opt,name=estado,proto3" json:"estado,omitempty"`
	IDTxc            uint64 `protobuf:"varint,6,opt,name=id_txc,proto3" json:"idTxc,omitempty"`
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
	IDImpuesto  int32  `json:"idImpuesto"`
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

// SmartContract Agrupador de funciones
type SmartContract struct {
	debug      bool
	isModeTest bool
}

//const startKey = "PER_20000000001"
//const endKey = "PER_35000000000"

var logger = shim.NewLogger("rut-afipcc")

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	t.debug = true
	return t.setInitImpuestos(stub)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	log.Print("=================================================================")
	log.Print("TxID [" + APIstub.GetTxID() + "]")
	function, args := APIstub.GetFunctionAndParameters()
	log.Print("Function [" + function + "] args [" + strings.Join(args, " ") + "]")

	switch function {
	case "putPersona":
		return s.putPersona(APIstub, args, JSON)
	case "putPersonaProto":
		return s.putPersona(APIstub, args, PROTOBUF)
	case "putPersonas":
		return s.putPersonas(APIstub, args, JSON)
	case "putPersonasProto":
		return s.putPersonas(APIstub, args, PROTOBUF)
	case "putParamImpuestos":
		return s.putParamImpuestos(APIstub, args)
	case "queryPersona":
		return s.queryPersona(APIstub, args)
	case "queryPersonasByRange":
		return s.queryPersonasByRange(APIstub, args)
	case "queryPersonaImpuestos":
		return s.queryPersonaImpuestos(APIstub, args)
	case "queryAllPersona":
		return s.queryAllPersona(APIstub)
	case "queryByKey":
		return s.queryByKey(APIstub, args)
	case "queryByKeyRange":
		return s.queryByKeyRange(APIstub, args)
	case "queryParamImpuestos":
		return s.queryParamImpuestos(APIstub)
	case "queryTxConfirmables":
		return s.queryTxConfirmables(APIstub, args)
	case "createTxConfirmable":
		return s.createTxConfirmable(APIstub, args)
	case "responseTxConfirmable":
		return s.responseTxConfirmable(APIstub, args)
	case "delPersona":
		return s.delPersona(APIstub, args)
	case "delPersonasByRange":
		return s.delPersonasByRange(APIstub, args)
	case "deleteAll":
		return s.deleteByKeyRange(APIstub, []string{"", ""})
	case "deleteByKeyRange":
		return s.deleteByKeyRange(APIstub, args)
	case "delParamImpuestosAll":
		return s.delParamImpuestosAll(APIstub)
	case "putPersonaImpuestos":
		return s.putPersonaImpuestos(APIstub, args)
	case "queryHistory":
		return s.queryHistory(APIstub, args)
	default:
		return s.clientErrorResponse(errors.New("Invalid Smart Contract function name " + function))
	}
}

func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("Ledger inicializado")
	return shim.Success(nil)
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
