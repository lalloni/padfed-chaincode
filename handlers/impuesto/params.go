package impuesto

import (
	"encoding/json"

	"github.com/pkg/errors"

	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/impuesto"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var (
	ImpuestoParam       = param.New("Impuesto JSON", parseImpuesto)
	ImpuestoListParam   = param.New("Impuesto List JSON", parseImpuestoList)
	CodigoImpuestoParam = param.SpecializeTyped(param.Uint64, "código impuesto", notZero)
)

var (
	impuestoListSchema = validator.MustLoadSchema("impuesto-list")
	impuestoSchema     = validator.MustLoadSchema("impuesto")
)

func parseImpuesto(bs []byte) (interface{}, error) {
	res, err := validator.Validate(impuestoSchema, bs)
	if err != nil {
		return nil, errors.Wrap(err, "validating impuesto")
	}
	if !res.Valid() {
		return nil, errors.Errorf("invalid impuesto: %s", res.String())
	}
	i := &model.Impuesto{}
	err = json.Unmarshal(bs, i)
	if err != nil {
		return nil, errors.Wrap(err, "parsing impuesto JSON")
	}
	return i, nil
}

func parseImpuestoList(bs []byte) (interface{}, error) {
	res, err := validator.Validate(impuestoListSchema, bs)
	if err != nil {
		return nil, errors.Wrap(err, "validating impuesto list")
	}
	if !res.Valid() {
		return nil, errors.Errorf("invalid impuesto list: %s", res.String())
	}
	is := []*model.Impuesto(nil)
	err = json.Unmarshal(bs, &is)
	if err != nil {
		return nil, errors.Wrap(err, "parsing impuesto list JSON")
	}
	r := []interface{}(nil)
	for _, i := range is {
		r = append(r, i)
	}
	return r, nil
}

func notZero(v interface{}) (interface{}, error) {
	if v.(uint64) == 0 {
		return nil, errors.New("must be greater than zero")
	}
	return v, nil
}
