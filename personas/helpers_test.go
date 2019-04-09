package personas_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/inscripciones"
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

func TestValimpuestosJSON(t *testing.T) {
	const cuit = 30679638943
	var per model.Persona
	err := json.Unmarshal(test.GetPersonaJSON(cuit), &per)
	if err != nil {
		t.Error("Error Failed to decode JSON of Impuestos")
	}
	if len(per.Impuestos) != 4 {
		t.Errorf("Persona debe tener 4 impuestos y tiene %d", len(per.Impuestos))
	}
	for _, imp := range per.Impuestos {
		v := inscripciones.GetImpuestoKeyByCuitID(cuit, imp.Impuesto)
		if v != "PER_"+strconv.FormatUint(cuit, 10)+"_IMP_"+strconv.FormatUint(uint64(imp.Impuesto), 10) {
			t.Errorf("1-Impuesto.Key no valido: %v", v)
		}
	}
}
