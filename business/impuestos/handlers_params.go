package impuestos

import (
	"encoding/json"

	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/pkg/errors"

	schemas "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/json-schemas"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var CodigoImpuestoParam = CodigoImpuestoParamVar(nil)

func CodigoImpuestoParamVar(ref *uint64) param.TypedParam {
	return param.SpecializeTyped(param.Uint64, "código impuesto", func(v interface{}) (interface{}, error) {
		c := v.(uint64)
		if c == 0 {
			return nil, errors.New("must be greater than zero")
		}
		if ref != nil {
			*ref = c
		}
		return c, nil
	})
}

var (
	ImpuestoParam     = param.New("Impuesto JSON", parseImpuesto)
	ImpuestoListParam = param.New("Impuesto List JSON", parseImpuestoList)
)

var (
	impuestoListSchema = schemas.MustLoadSchema("impuesto-list")
	impuestoSchema     = schemas.MustLoadSchema("impuesto")
)

func parseImpuesto(bs []byte) (interface{}, error) {
	res, err := validator.Validate(impuestoSchema, bs)
	if err != nil {
		return nil, errors.Wrap(err, "validating impuesto")
	}
	if !res.Valid() {
		return nil, errors.Errorf("invalid impuesto: %s", res.String())
	}
	i := &Impuesto{}
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
	is := []*Impuesto(nil)
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
