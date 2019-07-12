package personas

import (
	"encoding/json"
	"reflect"

	"github.com/lalloni/afip/cuit"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/pkg/errors"

	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var CUITParam = param.SpecializeTyped(param.Uint64, "CUIT", validateCUIT)

func validateCUIT(v interface{}) (interface{}, error) {
	if !cuit.IsValid(v.(uint64)) {
		return nil, errors.Errorf("invalid cuit/cuil number")
	}
	return v, nil
}

var PersonaParam = param.Typed("Persona JSON", reflect.TypeOf(&Persona{}), parsePersona)

func parsePersona(bs []byte) (interface{}, error) {

	res, err := validator.Validate(personaSchema, bs)
	if err != nil {
		return nil, errors.Wrap(err, "validating persona")
	}
	if !res.Valid() {
		return nil, errors.Errorf("invalid persona: %s", res.String())
	}

	per := &Persona{}
	if err = json.Unmarshal(bs, per); err != nil {
		return nil, errors.Wrap(err, "unmarshalling persona")
	}

	return per, nil

}

var personaSchema = validator.MustLoadSchema("persona")

var PersonaListParam = param.Typed("Persona List JSON", reflect.TypeOf([]Persona(nil)), parsePersonaList)

func parsePersonaList(bs []byte) (interface{}, error) {

	res, err := validator.Validate(personaListSchema, bs)
	if err != nil {
		return nil, errors.Wrap(err, "validating persona list")
	}
	if !res.Valid() {
		return nil, errors.Errorf("invalid persona list: %s", res.String())
	}

	pers := []*Persona{}
	err = json.Unmarshal(bs, &pers)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling persona list")
	}

	vs := []interface{}{}
	for _, per := range pers {
		vs = append(vs, per)
	}

	return vs, nil

}

var personaListSchema = validator.MustLoadSchema("persona-list")
