package personas_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/inscripciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/test"
)

func TestValidPersona(t *testing.T) {
	var persona model.Persona
	var personaJSON = test.GetPersonaJSON(30679638943)
	if err := personas.ArgToPersona([]byte(personaJSON), &persona); !err.IsOK() {
		t.Error(err.Msg)
	}
	if personas.GetPersonaKey(&persona) != "PER_30679638943" {
		t.Error("Persona.Key no valida " + personas.GetPersonaKey(&persona))
	}
	if err := personas.ArgToPersona([]byte("{error-dummy"), &persona); err.IsOK() {
		t.Error("JSON invalido, debe dar error " + err.Msg)
	}
}

func TestValimpuestosJSON(t *testing.T) {
	const cuit = 30679638943
	var impuestos model.Impuestos
	var personaJSON = test.GetPersonaJSON(cuit)
	err := json.Unmarshal([]byte(personaJSON), &impuestos)

	if err != nil {
		t.Error("Error Failed to decode JSON of Impuestos")
	}

	if len(impuestos.Impuestos) != 4 {
		t.Error("Persona debe tener 4 impuestos y tiene " + strconv.Itoa(len(impuestos.Impuestos)))
	}
	if inscripciones.GetImpuestoKeyByCuitID(cuit, impuestos.Impuestos[0].Impuesto) != "PER_30679638943_IMP_30" {
		t.Error("1-Impuesto.Key no valido " + inscripciones.GetImpuestoKeyByCuitID(cuit, impuestos.Impuestos[0].Impuesto))
	}
	if inscripciones.GetImpuestoKeyByCuitID(cuit, impuestos.Impuestos[3].Impuesto) != "PER_30679638943_IMP_34" {
		t.Error("3-Impuesto.Key no valido " + inscripciones.GetImpuestoKeyByCuitID(cuit, impuestos.Impuestos[3].Impuesto))
	}
}
