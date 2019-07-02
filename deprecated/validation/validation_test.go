package validation_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/test"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/validation"
	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/persona"
)

func TestValidPersona(t *testing.T) {
	var persona model.Persona
	var personaJSON = test.GetPersonaJSON(30679638943)
	if err := validation.ArgToPersona(personaJSON, &persona); !err.IsOK() {
		t.Error(err.Msg)
	}
	if err := validation.ArgToPersona([]byte("{error-dummy"), &persona); err.IsOK() {
		t.Error("JSON invalido, debe dar error " + err.Msg)
	}
}
