package ng

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/handlers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode/ng/support"
)

func Handlers() support.HandlerMap {
	return support.HandlerMap{
		"getPersona": handlers.GetPersonaHandler,
	}
}
