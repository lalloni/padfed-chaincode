package router

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
)

type Route struct {
	Name    string
	Check   authorization.Check
	Handler handler.Handler
}

type Config struct {
	Init *Route
	Funs []Route
}
