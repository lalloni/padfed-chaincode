package main

import (
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/chaincode"
)

// main function starts up the chaincode in the container during instantiate
func main() {
	log := shim.NewLogger("padfedcc")
	l := os.Getenv("SHIM_LOGGING_LEVEL")
	if l != "" {
		level, err := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
		if err != nil {
			log.Errorf("parsing SHIM_LOGGING_LEVEL with value %q: %v", l, err)
			os.Exit(1)
		}
		shim.SetLoggingLevel(level)
	}
	cc := chaincode.New(log, Version, false)
	if err := shim.Start(cc); err != nil {
		log.Errorf("starting chaincode: %v", err)
	}
}
