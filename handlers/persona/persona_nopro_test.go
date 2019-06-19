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

func TestDelPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("getp", nil, persona.GetPersonaHandler)
	r.SetHandler("putp", nil, persona.PutPersonaHandler)
	r.SetHandler("delpr", nil, persona.DelPersonaRangeHandler)

	mock := test.NewMock("test", r)

	pers := mtest.RandomPersonas(100, nil)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	min, max, index, ids := mtest.SummaryPersonasID(pers)

	_, res, payload, err := test.MockInvoke(t, mock, "delpr", min+1, max-1)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rids := []uint64{}
	err = mapstructure.Decode(payload.Content, &rids)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rids))
	a.ElementsMatch(ids[1:len(ids)-1], rids)

	_, res, payload, err = test.MockInvoke(t, mock, "getp", min)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	per := model.Persona{}
	err = mapstructure.Decode(payload.Content, &per)
	a.NoError(err)
	a.EqualValues(index[min], per)

	_, res, payload, err = test.MockInvoke(t, mock, "getp", max)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	per = model.Persona{}
	err = mapstructure.Decode(payload.Content, &per)
	a.NoError(err)
	a.EqualValues(index[max], per)

	for _, id := range ids[1 : len(ids)-1] {
		_, res, _, err = test.MockInvoke(t, mock, "getp", id)
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)
	}

}

func TestGetPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("putp", nil, persona.PutPersonaHandler)
	r.SetHandler("getpr", nil, persona.GetPersonaRangeHandler)

	mock := test.NewMock("test", r)

	pers := mtest.RandomPersonas(100, nil)
	min, max, index, ids := mtest.SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	_, res, payload, err := test.MockInvoke(t, mock, "getpr", min+1, max-1)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rpers))
	_, _, rindex, rids := mtest.SummaryPersonasID(rpers)
	_, ok := rindex[min]
	a.False(ok)
	_, ok = rindex[max]
	a.False(ok)
	for _, id := range rids {
		a.EqualValues(index[id], rindex[id])
	}

	_, res, payload, err = test.MockInvoke(t, mock, "getpr", 0, 99999999999)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers = []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	rmin, rmax, _, rids := mtest.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.EqualValues(min, rmin)
	a.EqualValues(max, rmax)
	a.ElementsMatch(pers, rpers)

}

func TestGetPersonaAllHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("putp", nil, persona.PutPersonaHandler)
	r.SetHandler("getpa", nil, persona.GetPersonaAllHandler)

	mock := test.NewMock("test", r)

	pers := mtest.RandomPersonas(100, nil)
	min, max, _, ids := mtest.SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "putp", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	_, res, payload, err := test.MockInvoke(t, mock, "getpa")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	rmin, rmax, _, rids := mtest.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.EqualValues(min, rmin)
	a.EqualValues(max, rmax)
	a.ElementsMatch(pers, rpers)

}
