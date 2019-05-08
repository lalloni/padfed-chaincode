package main

import (
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/generic"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func main() {

	log := shim.NewLogger("padfed")
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
	Everyone := authorization.Allowed

	r := router.New()

	r.SetInitHandler(OnlyAFIP, handler.SuccessHandler)

	// Meta
	r.SetHandlerFunc(Everyone, VersionHandler)

	// Business
	r.SetHandlerFunc(Everyone, persona.GetPersonaHandler)
	r.SetHandlerFunc(OnlyAFIP, persona.PutPersonaHandler)
	r.SetHandlerFunc(OnlyAFIP, persona.PutPersonaListHandler)
	r.SetHandlerFunc(OnlyAFIP, persona.DelPersonaHandler)
	r.SetHandlerFunc(OnlyAFIP, persona.DelPersonaRangeHandler)

	// Generic
	r.SetHandlerFunc(OnlyAFIP, generic.GetStateHandler)
	r.SetHandlerFunc(OnlyAFIP, generic.PutStateHandler)

	cc := chaincode.New(log, r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}

func VersionHandler(ctx *context.Context) *response.Response {
	return response.OK(Version)
}
