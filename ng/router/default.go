package router

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
)

func NameDefault(n string, h handler.Handler) Name {
	if n != "" {
		return Name(n)
	}
	return Name(handler.Name(h))
}

func CheckDefault(c, d authorization.Check) authorization.Check {
	if c == nil {
		return d
	}
	return c
}

func HandlerDefault(h, d handler.Handler) handler.Handler {
	if h == nil {
		return d
	}
	return h
}
