package main

import (
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
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
	All := authorization.Free

	rtr := router.New()

	rtr.SetInitHandler(OnlyAFIP, handler.SuccessHandler)
	rtr.SetInvokeHandler(All, "GetPersona", business.GetPersonaHandler)
	rtr.SetInvokeHandler(OnlyAFIP, "PutPersona", business.PutPersonaHandler)
	rtr.SetInvokeHandler(OnlyAFIP, "PutPersonaList", business.PutPersonaListHandler)
	rtr.SetInvokeHandler(OnlyAFIP, "DelPersona", business.DelPersonaHandler)
	rtr.SetInvokeHandler(OnlyAFIP, "DelPersonaRange", business.DelPersonaRangeHandler)

	cc := chaincode.New(log, rtr)

	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}

}
