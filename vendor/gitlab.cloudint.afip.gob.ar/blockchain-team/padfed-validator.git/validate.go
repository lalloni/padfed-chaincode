package validator

import (
	"github.com/lalloni/gojsonschema"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/jsoncheck"
)

func ValidateJSON(schema *gojsonschema.Schema, json []byte) (*ValidationResult, error) {
	err := jsoncheck.Check(json)
	if err != nil {
		return nil, errors.Wrap(err, "checking JSON syntax")
	}
	res, err := schema.Validate(gojsonschema.NewBytesLoader(json))
	if err != nil {
		return nil, errors.Wrap(err, "validating JSON document")
	}
	vr := ValidationResult{}
	for _, e := range res.Errors() {
		vr.Errors = append(vr.Errors, ValidationError{
			Field:       e.Field(),
			Description: e.Description(),
		})
	}
	return &vr, nil
}
