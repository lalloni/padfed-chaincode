package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/personas"
)

const AFIP = "AFIP"

type Handler func(shim.ChaincodeStubInterface, []string) *fabric.Response

type Handlers map[string]Handler

func BuildHandlers(testing bool) Handlers {

	h := Handlers{}

	// personas
	h["putPersona"] = requireMSP(AFIP, testing, personas.PutPersona)
	h["putPersonas"] = requireMSP(AFIP, testing, personas.PutPersonas)
	h["delPersona"] = requireMSP(AFIP, testing, personas.DelPersona)
	h["delPersonasByRange"] = requireMSP(AFIP, testing, personas.DelPersonasByRange)

	h["getPersona"] = personas.GetPersonaAPI
	h["queryPersona"] = personas.QueryPersona
	h["queryAllPersona"] = personas.QueryAllPersona

	// impuestos
	h["putParamImpuestos"] = impuestos.PutParamImpuestos
	h["queryParamImpuestos"] = impuestos.QueryParamImpuestos
	h["delParamImpuestosAll"] = adaptNoArg(impuestos.DeleteAll)

	// genéricas
	h["deleteAll"] = requireMSP(AFIP, testing, adaptNoArg(fabric.DeleteAll))
	h["deleteByKeyRange"] = requireMSP(AFIP, testing, adaptString2(fabric.DeleteByKeyRange))

	h["queryHistory"] = fabric.QueryHistory
	h["queryByKey"] = adaptString1(fabric.QueryByKey)
	h["queryByKeyRange"] = adaptString2(fabric.QueryByKeyRange)

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

func requireMSP(mspid string, testing bool, h Handler) Handler {
	return func(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
		if !testing {
			id, err := cid.New(stub)
			if err != nil {
				return fabric.SystemErrorResponse(fmt.Sprintf("getting client identity: %v", err))
			}
			s, err := id.GetMSPID()
			if err != nil {
				return fabric.SystemErrorResponse(fmt.Sprintf("getting client MSPID: %v", err))
			}
			if s != mspid {
				return fabric.ForbiddenErrorResponse(fmt.Sprintf("MSPID must be %q", mspid))
			}
		}
		return h(stub, args)
	}
}
