package store

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/filtering"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/marshaling"
)

type Option func(*simplestore)

func SetSep(sep *key.Sep) Option {
	return func(s *simplestore) {
		s.sep = sep
	}
}

func SetMarshaling(m marshaling.Marshaling) Option {
	return func(s *simplestore) {
		s.marshaling = m
	}
}

func SetFiltering(f filtering.Filtering) Option {
	return func(s *simplestore) {
		s.filtering = f
	}
}

func SetErrors(b bool) Option {
	return func(s *simplestore) {
		s.seterrs = b
	}
}
