package authorization

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

func Allowed(*context.Context) error {
	return nil
}

func Forbidden(*context.Context) error {
	return errors.New("not allowed")
}
