package validator

import (
	"strconv"
	"strings"

	"github.com/lalloni/gojsonschema"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/schemas"
)

type Validator interface {
	ValidatePersonaJSON(bs []byte) (*ValidationResult, error)
	ValidatePersonaListJSON(bs []byte) (*ValidationResult, error)
}

func New() (Validator, error) {
	personaSchema, err := schemas.PersonaSchema()
	if err != nil {
		return nil, errors.Wrap(err, "getting persona schema")
	}
	personaListSchema, err := schemas.PersonaListSchema()
	if err != nil {
		return nil, errors.Wrap(err, "getting persona schema")
	}
	return &validator{
		personaSchema:     personaSchema,
		personaListSchema: personaListSchema,
	}, nil
}

type validator struct {
	personaSchema     *gojsonschema.Schema
	personaListSchema *gojsonschema.Schema
}

func (v *validator) ValidatePersonaJSON(bs []byte) (*ValidationResult, error) {
	return ValidateJSON(v.personaSchema, bs)
}

func (v *validator) ValidatePersonaListJSON(bs []byte) (*ValidationResult, error) {
	return ValidateJSON(v.personaListSchema, bs)
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) String() string {
	l := len(r.Errors)
	sb := &strings.Builder{}
	sb.WriteString(strconv.Itoa(l))
	sb.WriteString(" error")
	if l != 1 {
		sb.WriteString("s")
	}
	for i, e := range r.Errors {
		if i == 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(e.Field)
		sb.WriteString(" ")
		sb.WriteString(e.Description)
		if i < l-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

type ValidationError struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}
