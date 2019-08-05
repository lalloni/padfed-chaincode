package personas

import (
	"encoding/json"
	"reflect"

	"github.com/lalloni/afip/cuit"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/pkg/errors"

	schemas "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/json-schemas"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var CUITParam = CUITParamVar(nil)

func CUITParamVar(ref *uint64) param.TypedParam {
	return param.SpecializeTyped(param.Uint64, "CUIT", func(v interface{}) (interface{}, error) {
		c := v.(uint64)
		if !cuit.IsValid(c) {
			return nil, errors.Errorf("invalid cuit/cuil number")
		}
		if ref != nil {
			*ref = c
		}
		return c, nil
	})
}

var PersonaParam = PersonaParamVar(nil)

func PersonaParamVar(ref *Persona) param.TypedParam {
	return param.Typed("Persona JSON", reflect.TypeOf(&Persona{}), func(bs []byte) (interface{}, error) {
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

		if ref != nil {
			*ref = *per
		}
		return per, nil
	})
}

var personaSchema = schemas.MustLoadSchema("persona")

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

var personaListSchema = schemas.MustLoadSchema("persona-list")

var EstadoParam = EstadoParamVar(nil)

func EstadoParamVar(ref *string) param.TypedParam {
	return param.SpecializeTyped(param.String, "Estado", func(v interface{}) (interface{}, error) {
		s := v.(string)
		if !estados[s] {
			return nil, errors.Errorf("invalid estado: '%s'", s)
		}
		if ref != nil {
			*ref = s
		}
		return s, nil
	})
}

var estados map[string]bool

func init() {
	valid := []string{"AC", "NA", "BD", "BP", "EX"}
	estados = make(map[string]bool, len(valid))
	for _, v := range valid {
		estados[v] = true
	}
}
