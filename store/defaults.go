package store

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/filtering"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/marshaling"
)

var (

	// DefaultMarshaling es el marshaling por defecto
	DefaultMarshaling = marshaling.JSON()

	// DefaultFiltering es el filtering por defecto
	DefaultFiltering = filtering.Copy()
)
