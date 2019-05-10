package persona_test

import (
	"net/http"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	r "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestGetPersonaHandler(t *testing.T) {
	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	per := &model.Persona{
		ID: 20242643772,
		Persona: &model.PersonaBasica{
			ID:       20242643772,
			Tipo:     "F",
			TipoID:   "C",
			Estado:   "A",
			Apellido: "Perez",
		},
	}

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("getp", nil, persona.GetPersonaHandler),
		r.R("putp", nil, persona.PutPersonaHandler))))

	res, payload, err := test.MockInvoke(mock, "getp", uint64(20242643772))
	t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.NotNil(payload)
	a.NotNil(res)
	a.EqualValues(http.StatusNotFound, res.Status)

	res, payload, err = test.MockInvoke(mock, "putp", per)
	t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.NotNil(payload)
	a.NotNil(res)
	a.EqualValues(http.StatusOK, res.Status)

	res, payload, err = test.MockInvoke(mock, "getp", uint64(20242643772))
	t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.NotNil(payload)
	a.NotNil(res)
	a.EqualValues(http.StatusOK, res.Status)

}
