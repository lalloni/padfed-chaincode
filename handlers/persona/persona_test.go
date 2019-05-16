package persona_test

import (
	"net/http"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	r "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestGetPutDelPersonaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("getp", nil, persona.GetPersonaHandler),
		r.R("putp", nil, persona.PutPersonaHandler),
		r.R("delp", nil, persona.DelPersonaHandler),
	)))

	for _, per := range test.RandomPersonas(10, nil) {
		per := per

		res, _, err := test.MockInvoke(mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(http.StatusNotFound, res.Status)

		res, _, err = test.MockInvoke(mock, "putp", per)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)

		res, payload, err := test.MockInvoke(mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		per1 := &model.Persona{}
		err = mapstructure.Decode(payload.Result, per1)
		a.NoError(err)
		a.NotNil(payload.Result)
		a.EqualValues(http.StatusOK, res.Status)
		a.EqualValues(&per, per1)

		res, _, err = test.MockInvoke(mock, "delp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)

		res, _, err = test.MockInvoke(mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(http.StatusNotFound, res.Status)

	}

}

func TestPutPersonaListHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	q := 50 // cantidad de personas a incluir en lista

	pers := test.RandomPersonas(q, nil)

	min, max, pi, _ := test.SummaryPersonasID(pers)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("putpl", nil, persona.PutPersonaListHandler),
		r.R("getpr", nil, persona.GetPersonaRangeHandler),
	)))

	res, payload, err := test.MockInvoke(mock, "putpl", &pers)
	a.NoError(err)
	a.NotNil(payload.Result)
	if !a.EqualValues(http.StatusOK, res.Status) {
		t.Logf("status: %d message: %q fault: %s list: %s", res.Status, res.Message, test.MustMarshal(payload.Fault), test.MustMarshal(pers))
	}
	a.EqualValues(q, payload.Result)

	res, payload, err = test.MockInvoke(mock, "getpr", min, max)
	a.NoError(err)
	a.NotNil(payload.Result)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(q, len(payload.Result.([]interface{})))
	for _, per := range payload.Result.([]interface{}) {
		per1 := model.Persona{}
		err = mapstructure.Decode(per, &per1)
		a.NoError(err)
		a.EqualValues(pi[per1.ID], per1)
	}

	res, payload, err = test.MockInvoke(mock, "getpr", min, max-1)
	a.NoError(err)
	a.NotNil(payload.Result)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(q-1, len(payload.Result.([]interface{})))

	res, payload, err = test.MockInvoke(mock, "getpr", min+1, max-1)
	a.NoError(err)
	a.NotNil(payload.Result)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(q-2, len(payload.Result.([]interface{})))

}
