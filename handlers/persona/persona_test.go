package persona_test

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	mtest "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/test"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response/status"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestGetPutDelPersonaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("getp", nil, persona.GetPersonaHandler)
	r.SetHandler("delp", nil, persona.DelPersonaHandler)
	r.SetHandler("putp", nil, persona.PutPersonaHandler)

	mock := test.NewMock("test", r)

	for _, per := range mtest.RandomPersonas(10, nil) {
		per := per

		_, res, _, err := test.MockInvoke(t, mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)

		_, res, _, err = test.MockInvoke(t, mock, "putp", per)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)

		_, res, payload, err := test.MockInvoke(t, mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		per1 := &model.Persona{}
		err = mapstructure.Decode(payload.Content, per1)
		a.NoError(err)
		a.NotNil(payload.Content)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues(&per, per1)

		_, res, _, err = test.MockInvoke(t, mock, "delp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)

		_, res, _, err = test.MockInvoke(t, mock, "getp", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)

	}

}

func TestPutPersonaListHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	q := 50 // cantidad de personas a incluir en lista

	pers := mtest.RandomPersonas(q, nil)

	min, max, pi, _ := mtest.SummaryPersonasID(pers)

	r := router.New()
	r.SetHandler("putpl", nil, persona.PutPersonaListHandler)
	r.SetHandler("getpr", nil, persona.GetPersonaRangeHandler)

	mock := test.NewMock("test", r)

	_, res, payload, err := test.MockInvoke(t, mock, "putpl", &pers)
	a.NoError(err)
	a.NotNil(payload.Content)
	if !a.EqualValues(status.OK, res.Status) {
		t.Logf("status: %d message: %q fault: %s list: %s", res.Status, res.Message, test.MustMarshal(payload.Fault), test.MustMarshal(pers))
	}
	a.EqualValues(q, payload.Content)

	_, res, payload, err = test.MockInvoke(t, mock, "getpr", min, max)
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(q, len(payload.Content.([]interface{})))
	for _, per := range payload.Content.([]interface{}) {
		per1 := model.Persona{}
		err = mapstructure.Decode(per, &per1)
		a.NoError(err)
		a.EqualValues(pi[per1.ID], per1)
	}

	_, res, payload, err = test.MockInvoke(t, mock, "getpr", min, max-1)
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(q-1, len(payload.Content.([]interface{})))

	_, res, payload, err = test.MockInvoke(t, mock, "getpr", min+1, max-1)
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(q-2, len(payload.Content.([]interface{})))

}

func TestGetPersonaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("getp", nil, persona.GetPersonaHandler)
	r.SetHandler("delp", nil, persona.DelPersonaHandler)

	mock := test.NewMock("test", r)

	for _, fun := range []string{"getp", "delp"} {
		_, res, _, err := test.MockInvoke(t, mock, fun, "foo", "bar")
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("1 argument expected (received 2)", res.Message)

		_, res, _, err = test.MockInvoke(t, mock, fun)
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("1 argument expected (received 0)", res.Message)

		_, res, _, err = test.MockInvoke(t, mock, fun, "-1")
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("CUIT argument 1: invalid natural integer: invalid syntax: '-1'", res.Message)
	}

}
