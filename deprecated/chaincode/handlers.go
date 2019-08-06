package chaincode

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/personas"
)

const AFIP = "AFIP"

type Handler func(shim.ChaincodeStubInterface, []string) *fabric.Response

type Handlers map[string]Handler

func BuildHandlers(version string, testing bool) Handlers {

	h := Handlers{}

	h["version"] = func(shim.ChaincodeStubInterface, []string) *fabric.Response {
		return fabric.SuccessResponseWithBuffer(bytes.NewBufferString(version))
	}

	// API Personas
	h["putPersona"] = onlyAFIP(testing, personas.PutPersona)
	h["getPersona"] = personas.GetPersonaAPI
	h["putPersonas"] = onlyAFIP(testing, personas.PutPersonas)

	// API Bajo Nivel
	h["queryPersona"] = personas.QueryPersona
	h["queryHistory"] = fabric.QueryHistory
	h["queryByKey"] = fabric.QueryByKey
	h["queryByKeyRange"] = fabric.QueryByKeyRange

	return h

}

func onlyAFIP(testing bool, h Handler) Handler {
	mspid := AFIP
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
