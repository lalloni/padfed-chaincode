package handler

import (
	"strings"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
)

func CheckArgsCount(ctx *context.Context, expected int) error {
	count := len(ctx.Stub.GetArgs()) - 1 // discount fn name in args[0]
	if expected != count {
		return errors.Errorf("argument count mismatch: received %d while expecting %d", count, expected)
	}
	return nil
}

func ExtractArgs(args [][]byte, pars ...param.Param) ([]interface{}, error) {
	parc := len(pars)
	argc := len(args)
	if parc != argc {
		return nil, errors.Errorf("argument count mismatch: received %d while expecting %d%s", argc, parc, names(pars))
	}
	res := []interface{}(nil)
	for i, par := range pars {
		v, err := par.From(args[i])
		if err != nil {
			return nil, errors.Wrapf(err, "%s argument %d", par.Name(), i+1)
		}
		res = append(res, v)
	}
	return res, nil
}

func names(pars []param.Param) string {
	ss := []string(nil)
	for _, par := range pars {
		ss = append(ss, par.Name())
	}
	s := strings.Join(ss, ", ")
	if len(s) > 0 {
		s = " (" + s + ")"
	}
	return s
}
