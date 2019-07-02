package persona

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/afip/cuit"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/persona"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response/status"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestDelPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	pers := model.RandomPersonas(10, nil)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "PutPersona", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	index, ids := model.SummaryPersonasID(pers)
	min := ids[0]
	max := ids[len(ids)-1]

	_, res, payload, err := test.MockInvoke(t, mock, "DelPersonaRange", cuit.Succ(min), cuit.Pred(max))
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rids := []uint64{}
	err = mapstructure.Decode(payload.Content, &rids)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rids))
	a.ElementsMatch(ids[1:len(ids)-1], rids)

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersona", min)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	per := model.Persona{}
	err = mapstructure.Decode(payload.Content, &per)
	a.NoError(err)
	a.EqualValues(index[min], per)

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersona", max)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	per = model.Persona{}
	err = mapstructure.Decode(payload.Content, &per)
	a.NoError(err)
	a.EqualValues(index[ids[len(ids)-1]], per)

	for _, id := range ids[1 : len(ids)-1] {
		_, res, _, err = test.MockInvoke(t, mock, "GetPersona", id)
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)
	}

}

func TestGetPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	pers := model.RandomPersonas(100, nil)
	index, ids := model.SummaryPersonasID(pers)
	min := ids[0]
	max := ids[len(ids)-1]

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "PutPersona", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaRange", cuit.Succ(min), cuit.Pred(max))
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rpers))
	rindex, rids := model.SummaryPersonasID(rpers)
	_, ok := rindex[min]
	a.False(ok)
	_, ok = rindex[max]
	a.False(ok)
	for _, id := range rids {
		a.EqualValues(index[id], rindex[id])
	}

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersonaRange", cuit.Min, cuit.Max)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers = []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	_, rids = model.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.ElementsMatch(pers, rpers)

}

func TestGetPersonaAllHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	pers := model.RandomPersonas(100, nil)
	_, ids := model.SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "PutPersona", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaAll")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers := []model.Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	_, rids := model.SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.ElementsMatch(pers, rpers)

}
