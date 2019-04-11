package personas_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestGetPersona(t *testing.T) {
	stub := test.SetInitTests(t)
	var cuit uint64 = 30679638943
	res := test.PutPersona(t, stub, cuit)
	if res.Status != shim.OK {
		t.Error("putPersona fail: " + res.Message)
	}
	res = test.GetPersona(t, stub, cuit)
	if res.Status != shim.OK {
		t.Error(res.Message)
	}
	fmt.Println("getPersona ", string(res.Payload))

	personaGet := model.Persona{}
	if err := json.Unmarshal(res.Payload, &personaGet); err != nil {
		t.Errorf("unmarshalling response message: %v", err)
	}

	var personaJSON = test.GetPersonaJSON(cuit)
	personaPut := model.Persona{}
	if err := json.Unmarshal(personaJSON, &personaPut); err != nil {
		t.Errorf("unmarshalling response message: %v", err)
	}

	if personaPut.ID != personaGet.ID {
		t.Errorf("los id son distintos")
	}

	if personaPut.Persona != personaPut.Persona {
		t.Errorf("las Personas son distintos")
	}

	if len(personaPut.Impuestos) != len(personaPut.Impuestos) {
		t.Errorf("los Impuestos son distintos ")
	}
}
