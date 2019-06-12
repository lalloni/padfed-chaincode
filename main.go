package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/generic"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/logging"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

const name = "padfedcc"

var (
	OnlyAFIP = authorization.MSPID("AFIP")
	Free     = authorization.Allowed
)

func main() {

	log := logging.Setup(name)

	r := router.New(nil)

	r.SetInitHandler(OnlyAFIP, handler.SuccessHandler)

	// Meta
	r.SetHandler("GetVersion", Free, VersionHandler)
	r.SetHandler("GetFunctions", Free, r.FunctionsHandler())

	// Business
	r.SetHandler("GetPersona", Free, persona.GetPersonaHandler)
	r.SetHandler("DelPersona", OnlyAFIP, persona.DelPersonaHandler)
	r.SetHandler("PutPersona", OnlyAFIP, persona.PutPersonaHandler)
	r.SetHandler("PutPersonaList", OnlyAFIP, persona.PutPersonaListHandler)

	// Business (debugging)
	r.SetHandler("GetPersonaRange", Free, persona.GetPersonaRangeHandler)
	r.SetHandler("DelPersonaRange", OnlyAFIP, persona.DelPersonaRangeHandler)
	r.SetHandler("GetPersonaAll", Free, persona.GetPersonaAllHandler)

	// Generic
	r.SetHandler("GetStates", Free, generic.GetStatesHandler)
	r.SetHandler("DelStates", OnlyAFIP, generic.DelStatesHandler)
	r.SetHandler("PutStates", OnlyAFIP, generic.PutStatesHandler)

	// History
	r.SetHandler("GetStatesHistory", Free, generic.GetStatesHistoryHandler)

	// ======================== Deprecated handlers ============================
	//
	// No se necesitan ya que su funcionalidad está duplicada de los de arriba y
	// no presentan al cliente una interfaz robusta, consistente y homogenea.
	//
	// TODO eliminar bloque antes de 1.0.0
	r.SetHandler("putPersona", OnlyAFIP, deprecated.Adapter(personas.PutPersona, "PutPersona"))
	r.SetHandler("delPersona", OnlyAFIP, deprecated.Adapter(personas.DelPersona, "DelPersona"))
	r.SetHandler("getPersona", Free, deprecated.Adapter(personas.GetPersonaAPI, "GetPersona"))
	r.SetHandler("putPersonas", OnlyAFIP, deprecated.Adapter(personas.PutPersonas, "PutPersonaList"))
	r.SetHandler("delPersonasByRange", OnlyAFIP, deprecated.Adapter(personas.DelPersonasByRange, "DelPersonaRange"))
	r.SetHandler("deleteAll", OnlyAFIP, deprecated.Adapter(fabric.DeleteAll, "DelStates"))
	r.SetHandler("deleteByKeyRange", OnlyAFIP, deprecated.Adapter(fabric.DeleteByKeyRange, "DelStates"))
	r.SetHandler("queryPersona", OnlyAFIP, deprecated.Adapter(personas.QueryPersona, "GetStates"))
	r.SetHandler("queryAllPersona", OnlyAFIP, deprecated.Adapter(personas.QueryAllPersona, "GetStates"))
	r.SetHandler("queryHistory", OnlyAFIP, deprecated.Adapter(fabric.QueryHistory, "GetStatesHistory"))
	r.SetHandler("queryByKey", OnlyAFIP, deprecated.Adapter(fabric.QueryByKey, "GetStates"))
	r.SetHandler("queryByKeyRange", OnlyAFIP, deprecated.Adapter(fabric.QueryByKeyRange, "GetStates"))
	//
	// =========================================================================

	cc := chaincode.New(name, Version, r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}

func VersionHandler(ctx *context.Context) *response.Response {
	return response.OK(Version)
}
