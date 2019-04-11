package personas_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestGetPersona(t *testing.T) {
	a := assert.New(t)

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

	a.Equal(personaPut.Persona, personaGet.Persona, "personas distintas")

	a.Equal(personaPut.Impuestos, personaGet.Impuestos, "impuestos distintas")

}
