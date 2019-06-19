package persona

import (
	"encoding/json"
	"reflect"

	"github.com/lalloni/afip/cuit"

	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var CUITParam = param.Uint64.Specialize("CUIT", validateCUIT)

func validateCUIT(v interface{}) (interface{}, error) {
	if !cuit.IsValid(v.(uint64)) {
		return nil, errors.Errorf("invalid cuit/cuil number")
	}
	return v, nil
}

var PersonaParam = param.New("Persona JSON", reflect.TypeOf(&model.Persona{}), parsePersona)

func parsePersona(args [][]byte) (interface{}, int, error) {

	bs := args[0]

	res, err := validator.ValidatePersonaJSON(bs)
	if err != nil {
		return nil, 0, errors.Wrap(err, "validating persona")
	}
	if !res.Valid() {
		return nil, 0, errors.Errorf("invalid persona: %s", res.String())
	}

	per := &model.Persona{}
	if err = json.Unmarshal(bs, per); err != nil {
		return nil, 0, errors.Wrap(err, "unmarshalling persona")
	}

	return per, 1, nil

}

var PersonaListParam = param.New("Persona List JSON", reflect.TypeOf([]model.Persona(nil)), parsePersonaList)

func parsePersonaList(args [][]byte) (interface{}, int, error) {

	bs := args[0]

	res, err := validator.ValidatePersonaListJSON(bs)
	if err != nil {
		return nil, 0, errors.Wrap(err, "validating persona list")
	}
	if !res.Valid() {
		return nil, 0, errors.Errorf("invalid persona list: %s", res.String())
	}

	pers := []model.Persona{}
	err = json.Unmarshal(bs, &pers)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unmarshalling persona list")
	}

	return pers, 1, nil

}
