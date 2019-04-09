package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/inscripciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/personas"
)

type Handler func(shim.ChaincodeStubInterface, []string) *fabric.Response

type Handlers map[string]Handler

func BuildHandlers() Handlers {

	h := Handlers{}

	// personas
	h["putPersona"] = personas.PutPersona
	h["putPersonas"] = personas.PutPersonas
	h["delPersona"] = personas.DelPersona
	h["delPersonasByRange"] = personas.DelPersonasByRange
	h["queryPersona"] = personas.QueryPersona
	h["queryPersonasByRange"] = personas.QueryPersonasByRange
	h["queryAllPersona"] = personas.QueryAllPersona

	// impuestos
	h["putParamImpuestos"] = impuestos.PutParamImpuestos
	h["queryParamImpuestos"] = impuestos.QueryParamImpuestos
	h["delParamImpuestosAll"] = adaptNoArg(impuestos.DeleteAll)

	// inscripciones
	h["putPersonaImpuestos"] = inscripciones.PutPersonaImpuestos
	h["queryPersonaImpuestos"] = inscripciones.QueryPersonaImpuestos

	// genéricas
	h["queryHistory"] = fabric.QueryHistory
	h["queryByKey"] = adaptString1(fabric.QueryByKey)
	h["queryByKeyRange"] = adaptString2(fabric.QueryByKeyRange)
	h["deleteAll"] = adaptNoArg(fabric.DeleteAll)
	h["deleteByKeyRange"] = adaptString2(fabric.DeleteByKeyRange)

	return h

}

func adaptNoArg(h func(shim.ChaincodeStubInterface) *fabric.Response) Handler {
	return func(s shim.ChaincodeStubInterface, _ []string) *fabric.Response {
		return h(s)
	}
}

func adaptString1(h func(shim.ChaincodeStubInterface, string) *fabric.Response) Handler {
	return func(s shim.ChaincodeStubInterface, args []string) *fabric.Response {
		return h(s, args[0])
	}
}

func adaptString2(h func(shim.ChaincodeStubInterface, string, string) *fabric.Response) Handler {
	return func(s shim.ChaincodeStubInterface, args []string) *fabric.Response {
		return h(s, args[0], args[1])
	}
}
