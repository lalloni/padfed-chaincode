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

	r := router.New(nil)

	r.SetInitHandler(OnlyAFIP, nil)

	// Meta
	r.SetHandler("version", nil, VersionHandler)

	// Business
	r.SetHandler("GetPersona", nil, persona.GetPersonaHandler)
	r.SetHandler("DelPersona", OnlyAFIP, persona.DelPersonaHandler)
	r.SetHandler("PutPersona", OnlyAFIP, persona.PutPersonaHandler)
	r.SetHandler("PutPersonaList", OnlyAFIP, persona.PutPersonaListHandler)

	// Business not productive
	r.SetHandler("GetPersonaRange", OnlyAFIP, persona.GetPersonaRangeHandler)
	r.SetHandler("DelPersonaRange", OnlyAFIP, persona.DelPersonaRangeHandler)

	// Generic
	r.SetHandler("GetState", OnlyAFIP, generic.GetStateHandler)
	r.SetHandler("PutState", OnlyAFIP, generic.PutStateHandler)

	cc := chaincode.New("padfed", r)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}

func VersionHandler(ctx *context.Context) *response.Response {
	return response.OK(Version)
}
