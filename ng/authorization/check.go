package authorization

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

type Check func(*context.Context) error
