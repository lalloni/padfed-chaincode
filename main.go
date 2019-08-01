package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/logging"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	deprecatedpersonas "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/state"
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
	personas.AddHandlers(r)
	impuestos.AddHandlers(r)

	// States
	state.AddHandlers(r)

	// ======================== Deprecated handlers ============================
	//
	// No se necesitan ya que su funcionalidad está duplicada de los de arriba y
	// no presentan al cliente una interfaz robusta, consistente y homogenea.
	//
	// TODO eliminar bloque antes de 1.0.0
	r.SetHandler("putPersona", common.AFIP, deprecated.WarningAdapter(deprecatedpersonas.PutPersona, "PutPersona"))
	r.SetHandler("delPersona", common.AFIP, deprecated.WarningAdapter(deprecatedpersonas.DelPersona, "DelPersona"))
	r.SetHandler("getPersona", common.Free, deprecated.WarningAdapter(deprecatedpersonas.GetPersonaAPI, "GetPersona"))
	r.SetHandler("putPersonas", common.AFIP, deprecated.WarningAdapter(deprecatedpersonas.PutPersonas, "PutPersonaList"))
	r.SetHandler("delPersonasByRange", common.AFIP, deprecated.WarningAdapter(deprecatedpersonas.DelPersonasByRange, "DelPersonaRange"))
	r.SetHandler("deleteAll", common.AFIP, deprecated.WarningAdapter(fabric.DeleteAll, "DelStates"))
	r.SetHandler("deleteByKeyRange", common.AFIP, deprecated.WarningAdapter(fabric.DeleteByKeyRange, "DelStates"))
	r.SetHandler("queryPersona", common.Free, deprecated.WarningAdapter(deprecatedpersonas.QueryPersona, "GetStates"))
	r.SetHandler("queryAllPersona", common.Free, deprecated.WarningAdapter(deprecatedpersonas.QueryAllPersona, "GetStates"))
	r.SetHandler("queryHistory", common.Free, deprecated.Adapter(fabric.QueryHistory))
	r.SetHandler("queryByKey", common.Free, deprecated.WarningAdapter(fabric.QueryByKey, "GetStates"))
	r.SetHandler("queryByKeyRange", common.Free, deprecated.WarningAdapter(fabric.QueryByKeyRange, "GetStates"))
	//
	// =========================================================================

	cc := chaincode.New(name, Version, r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}
