package common

import "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"

var (
	AFIP = authorization.MSPID("AFIP")
	Free = authorization.Allowed
)
