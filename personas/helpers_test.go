package personas_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestValidPersona(t *testing.T) {
	var persona model.Persona
	var personaJSON = test.GetPersonaJSON(30679638943)
	if err := model.ArgToPersona(personaJSON, &persona); !err.IsOK() {
		t.Error(err.Msg)
	}
	if personas.GetPersonaKey(&persona) != "PER_30679638943" {
		t.Error("Persona.Key no valida " + personas.GetPersonaKey(&persona))
	}
	if err := model.ArgToPersona([]byte("{error-dummy"), &persona); err.IsOK() {
		t.Error("JSON invalido, debe dar error " + err.Msg)
	}
}
