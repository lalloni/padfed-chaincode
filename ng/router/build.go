package router

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler"
)

func Rs(rs ...*Route) []Route {
	res := []Route{}
	for _, r := range rs {
		if r != nil {
			res = append(res, *r)
		}
	}
	return res
}

func R(n string, c authorization.Check, h handler.Handler) *Route {
	return &Route{
		Name:    n,
		Check:   c,
		Handler: h,
	}
}

func RCH(c authorization.Check, h handler.Handler) *Route {
	return &Route{
		Check:   c,
		Handler: h,
	}
}

func RH(h handler.Handler) *Route {
	return &Route{
		Handler: h,
	}
}

func C(init *Route, funs ...*Route) *Config {
	return &Config{
		Init: init,
		Funs: Rs(funs...),
	}
}
