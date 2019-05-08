package main

import (
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/generic"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

const name = "padfedcc"

func main() {

	log := shim.NewLogger(name)
	l := os.Getenv("SHIM_LOGGING_LEVEL")
	if l != "" {
		level, err := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
		if err != nil {
			log.Errorf("parsing SHIM_LOGGING_LEVEL from value %q: %v", l, err)
			os.Exit(1)
		}
		shim.SetLoggingLevel(level)
	}

	OnlyAFIP := authorization.MSPID("AFIP")

	r := router.New(nil)

	r.SetInitHandler(OnlyAFIP, nil)

	// Meta
	r.SetHandler("Version", nil, VersionHandler)
	r.SetHandler("Functions", nil, r.FunctionsHandler())

	// Business
	r.SetHandler("GetPersona", nil, persona.GetPersonaHandler)
	r.SetHandler("DelPersona", OnlyAFIP, persona.DelPersonaHandler)
	r.SetHandler("PutPersona", OnlyAFIP, persona.PutPersonaHandler)
	r.SetHandler("PutPersonaList", OnlyAFIP, persona.PutPersonaListHandler)

	// Business (debugging only)
	r.SetHandler("GetPersonaRange", OnlyAFIP, persona.GetPersonaRangeHandler)
	r.SetHandler("DelPersonaRange", OnlyAFIP, persona.DelPersonaRangeHandler)
	r.SetHandler("GetPersonaAll", OnlyAFIP, persona.GetPersonaAllHandler)

	// Generic
	r.SetHandler("GetStates", OnlyAFIP, generic.GetStatesHandler)
	r.SetHandler("PutStates", OnlyAFIP, generic.PutStatesHandler)
	r.SetHandler("DelStates", OnlyAFIP, generic.DelStatesHandler)

	cc := chaincode.New(name, r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}

func VersionHandler(ctx *context.Context) *response.Response {
	return response.OK(Version)
}
