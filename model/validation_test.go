package model_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestValidPersona(t *testing.T) {
	var persona model.Persona
	var personaJSON = test.GetPersonaJSON(30679638943)
	if err := model.ArgToPersona(personaJSON, &persona); !err.IsOK() {
		t.Error(err.Msg)
	}
	if err := model.ArgToPersona([]byte("{error-dummy"), &persona); err.IsOK() {
		t.Error("JSON invalido, debe dar error " + err.Msg)
	}
}
