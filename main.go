package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/logging"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/impuesto"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/state"
)

const name = "padfedcc"

var log = logging.Setup(name)

func main() {

	r := router.New()

	r.SetInitHandler(common.AFIP, handler.SuccessHandler)

	// Meta
	r.SetHandler("GetVersion", common.Free, handler.ValueHandler(Version))
	r.SetHandler("GetFunctions", common.Free, router.FunctionsHandler(r))

	// Business
	persona.AddHandlers(r)
	impuesto.AddHandlers(r)

	// States
	state.AddHandlers(r)

	// ======================== Deprecated handlers ============================
	//
	// No se necesitan ya que su funcionalidad está duplicada de los de arriba y
	// no presentan al cliente una interfaz robusta, consistente y homogenea.
	//
	// TODO eliminar bloque antes de 1.0.0
	r.SetHandler("putPersona", common.AFIP, deprecated.Adapter(personas.PutPersona, "PutPersona"))
	r.SetHandler("delPersona", common.AFIP, deprecated.Adapter(personas.DelPersona, "DelPersona"))
	r.SetHandler("getPersona", common.Free, deprecated.Adapter(personas.GetPersonaAPI, "GetPersona"))
	r.SetHandler("putPersonas", common.AFIP, deprecated.Adapter(personas.PutPersonas, "PutPersonaList"))
	r.SetHandler("delPersonasByRange", common.AFIP, deprecated.Adapter(personas.DelPersonasByRange, "DelPersonaRange"))
	r.SetHandler("deleteAll", common.AFIP, deprecated.Adapter(fabric.DeleteAll, "DelStates"))
	r.SetHandler("deleteByKeyRange", common.AFIP, deprecated.Adapter(fabric.DeleteByKeyRange, "DelStates"))
	r.SetHandler("queryPersona", common.Free, deprecated.Adapter(personas.QueryPersona, "GetStates"))
	r.SetHandler("queryAllPersona", common.Free, deprecated.Adapter(personas.QueryAllPersona, "GetStates"))
	r.SetHandler("queryHistory", common.Free, deprecated.Adapter(fabric.QueryHistory, "GetStatesHistory"))
	r.SetHandler("queryByKey", common.Free, deprecated.Adapter(fabric.QueryByKey, "GetStates"))
	r.SetHandler("queryByKeyRange", common.Free, deprecated.Adapter(fabric.QueryByKeyRange, "GetStates"))
	//
	// =========================================================================

	cc := chaincode.New(name, Version, r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}
