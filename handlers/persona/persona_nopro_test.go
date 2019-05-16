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

func TestDelPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("delpr", nil, persona.DelPersonaRangeHandler),
		r.R("putp", nil, persona.PutPersonaHandler),
		r.R("getp", nil, persona.GetPersonaHandler),
	)))

	pers := test.RandomPersonas(100, nil)

	for _, per := range pers {
		per := per
		res, _, err := test.MockInvoke(mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	}

	min, max, index, ids := test.SummaryPersonasID(pers)

	res, payload, err := test.MockInvoke(mock, "delpr", min+1, max-1)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	rids := []uint64{}
	err = mapstructure.Decode(payload.Result, &rids)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rids))
	a.ElementsMatch(ids[1:len(ids)-1], rids)

	res, payload, err = test.MockInvoke(mock, "getp", min)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	per := model.Persona{}
	err = mapstructure.Decode(payload.Result, &per)
	a.NoError(err)
	a.EqualValues(index[min], per)

	res, payload, err = test.MockInvoke(mock, "getp", max)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	per = model.Persona{}
	err = mapstructure.Decode(payload.Result, &per)
	a.NoError(err)
	a.EqualValues(index[max], per)

	for _, id := range ids[1 : len(ids)-1] {
		res, _, err = test.MockInvoke(mock, "getp", id)
		a.NoError(err)
		a.EqualValues(http.StatusNotFound, res.Status)
	}

}

func TestGetPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("putp", nil, persona.PutPersonaHandler),
		r.R("getpr", nil, persona.GetPersonaRangeHandler),
	)))

	pers := test.RandomPersonas(100, nil)
	min, max, index, ids := test.SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		res, _, err := test.MockInvoke(mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	}

	res, payload, err := test.MockInvoke(mock, "getpr", min+1, max-1)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Result, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rpers))
	_, _, rindex, rids := test.SummaryPersonasID(rpers)
	_, ok := rindex[min]
	a.False(ok)
	_, ok = rindex[max]
	a.False(ok)
	for _, id := range rids {
		a.EqualValues(index[id], rindex[id])
	}

	res, payload, err = test.MockInvoke(mock, "getpr", 0, 99999999999)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	rpers = []model.Persona{}
	err = mapstructure.Decode(payload.Result, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	rmin, rmax, _, rids := test.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.EqualValues(min, rmin)
	a.EqualValues(max, rmax)
	a.ElementsMatch(pers, rpers)

}

func TestGetPersonaAllHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("putp", nil, persona.PutPersonaHandler),
		r.R("getpa", nil, persona.GetPersonaAllHandler),
	)))

	pers := test.RandomPersonas(100, nil)
	min, max, _, ids := test.SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		res, _, err := test.MockInvoke(mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	}

	res, payload, err := test.MockInvoke(mock, "getpa")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Result, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	rmin, rmax, _, rids := test.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.EqualValues(min, rmin)
	a.EqualValues(max, rmax)
	a.ElementsMatch(pers, rpers)

}
