package handler

import (
	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

func ValidateArgCount(ctx *context.Context, expected int) error {
	count := len(ctx.Stub.GetArgs()) - 1 // discount function name in args[0]
	if expected != count {
		s := ""
		if expected > 1 {
			s = "s"
		}
		return errors.Errorf("%d argument%s expected (received %d)", expected, s, count)
	}
	return nil
}
