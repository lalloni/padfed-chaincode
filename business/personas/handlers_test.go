package personas

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/afip/cuit"
	"github.com/lalloni/fabrikit/chaincode/response/status"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestGetPutDelPersonaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	for _, per := range RandomPersonas(10, nil) {
		per := per

		_, res, _, err := test.MockInvoke(t, mock, "GetPersona", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)

		_, res, _, err = test.MockInvoke(t, mock, "PutPersona", per)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)

		_, res, payload, err := test.MockInvoke(t, mock, "GetPersona", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		per1 := &Persona{}
		err = mapstructure.Decode(payload.Content, per1)
		a.NoError(err)
		a.NotNil(payload.Content)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues(&per, per1)

		_, res, _, err = test.MockInvoke(t, mock, "DelPersona", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)

		_, res, _, err = test.MockInvoke(t, mock, "GetPersona", per.ID)
		t.Logf("response status: %v message: %s payload: %s", res.Status, res.Message, string(res.Payload))
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)

	}

}

func TestPutPersonaListHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	q := 50 // cantidad de personas a incluir en lista

	pers := RandomPersonas(q, nil)

	pi, pids := SummaryPersonasID(pers)
	min := pids[0]
	max := pids[len(pids)-1]

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	_, res, payload, err := test.MockInvoke(t, mock, "PutPersonaList", &pers)
	a.NoError(err)
	a.NotNil(payload.Content)
	if !a.EqualValues(status.OK, res.Status) {
		t.Logf("status: %d message: %q fault: %s list: %s", res.Status, res.Message, test.MustMarshal(payload.Fault), test.MustMarshal(pers))
	}
	a.EqualValues(q, payload.Content)

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersonaRange", min, max)
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(q, len(payload.Content.([]interface{})))
	for _, per := range payload.Content.([]interface{}) {
		per1 := Persona{}
		err = mapstructure.Decode(per, &per1)
		a.NoError(err)
		a.EqualValues(pi[per1.ID], per1)
	}

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersonaRange", min, cuit.Pred(max))
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	if payload.Content != nil {
		a.EqualValues(q-1, len(payload.Content.([]interface{})))
	}

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersonaRange", cuit.Succ(min), max)
	a.NoError(err)
	a.NotNil(payload.Content)
	a.EqualValues(status.OK, res.Status)
	if payload.Content != nil {
		a.EqualValues(q-1, len(payload.Content.([]interface{})))
	}
}

func TestGetPersonaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	for _, fun := range []string{"GetPersona", "DelPersona"} {
		_, res, _, err := test.MockInvoke(t, mock, fun, "foo", "bar")
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("invalid persona id: argument count mismatch: received 2 while expecting 1 (CUIT)", res.Message)

		_, res, _, err = test.MockInvoke(t, mock, fun)
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("invalid persona id: argument count mismatch: received 0 while expecting 1 (CUIT)", res.Message)

		_, res, _, err = test.MockInvoke(t, mock, fun, "-1")
		a.NoError(err)
		a.EqualValues(400, res.Status)
		a.EqualValues("invalid persona id: CUIT argument 1: invalid natural integer: invalid syntax: '-1'", res.Message)
	}

}

func TestDelPersonaRangeHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	pers := RandomPersonas(10, nil)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "PutPersona", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	index, ids := SummaryPersonasID(pers)
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
	per := Persona{}
	err = mapstructure.Decode(payload.Content, &per)
	a.NoError(err)
	a.EqualValues(index[min], per)

	_, res, payload, err = test.MockInvoke(t, mock, "GetPersona", max)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	per = Persona{}
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

	pers := RandomPersonas(100, nil)
	index, ids := SummaryPersonasID(pers)
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
	rpers := []Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers)-2, len(rpers))
	rindex, rids := SummaryPersonasID(rpers)
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
	rpers = []Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	_, rids = SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.ElementsMatch(pers, rpers)

}

func TestGetPersonaAllHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	pers := RandomPersonas(100, nil)
	_, ids := SummaryPersonasID(pers)

	for _, per := range pers {
		per := per
		_, res, _, err := test.MockInvoke(t, mock, "PutPersona", &per)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	}

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaAll")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	rpers := []Persona{}
	err = mapstructure.Decode(payload.Content, &rpers)
	a.NoError(err)
	a.EqualValues(len(pers), len(rpers))
	_, rids := SummaryPersonasID(rpers)
	a.EqualValues(ids, rids)
	a.ElementsMatch(pers, rpers)
}

func TestGetPersonaBasicaHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	per := RandomPersonas(1, nil)[0]

	_, res, _, err := test.MockInvoke(t, mock, "PutPersona", per)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaBasica", per.ID)
	a.NoError(err)
	a.NotNil(payload)
	a.EqualValues(status.OK, res.Status)
	var pb Basica
	a.NoError(mapstructure.Decode(payload.Content, &pb))
	a.EqualValues(per.Persona, &pb)

}

func TestGetPersonaCollectionHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	per := RandomPersonas(1, nil)[0]

	_, res, _, err := test.MockInvoke(t, mock, "PutPersona", per)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)

	for _, col := range Schema.Collections() {
		col := col
		t.Run(col.Name, func(t *testing.T) {
			_, res, payload, err := test.MockInvoke(t, mock, "GetPersona"+col.Name, per.ID)
			a.NoError(err)
			a.NotNil(payload)
			colval := reflect.ValueOf(col.Getter(per))
			if colval.Len() == 0 {
				a.EqualValues(status.NotFound, res.Status)
				return
			}
			a.EqualValues(status.OK, res.Status)
			bs, err := json.Marshal(col.Getter(per))
			a.NoError(err)
			want := map[string]interface{}{}
			a.NoError(json.Unmarshal(bs, &want))
			a.EqualValues(want, payload.Content)
		})
	}

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaEtiquetas", per.ID)
	a.NoError(err)
	a.NotNil(payload)
	a.EqualValues(status.OK, res.Status)
	var es map[string]*Etiqueta
	a.NoError(mapstructure.Decode(payload.Content, &es))
	a.EqualValues(per.Etiquetas, es)

}

func TestGetPersonaCollectionItemHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	addTestingHandlers(r)

	mock := test.NewMock("test", r)

	per := RandomPersonas(1, nil)[0]

	_, res, _, err := test.MockInvoke(t, mock, "PutPersona", per)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)

	for _, col := range Schema.Collections() {
		col := col
		t.Run(col.Name, func(t *testing.T) {
			colval := reflect.ValueOf(col.Getter(per))
			if colval.Len() == 0 {
				t.Skipf("collection %s is empty", col.Name)
				return
			}
			item := colval.MapKeys()[0].String()
			_, res, payload, err := test.MockInvoke(t, mock, "GetPersona"+col.Name+"Item", per.ID, item)
			a.NoError(err)
			a.NotNil(payload)
			a.EqualValues(status.OK, res.Status)
			itemval := col.ItemCreator()
			a.NoError(mapstructure.Decode(payload.Content, &itemval))
			a.EqualValues(colval.MapIndex(reflect.ValueOf(item)).Interface(), itemval)
		})
	}

	_, res, payload, err := test.MockInvoke(t, mock, "GetPersonaEtiquetas", per.ID)
	a.NoError(err)
	a.NotNil(payload)
	a.EqualValues(status.OK, res.Status)
	var es map[string]*Etiqueta
	a.NoError(mapstructure.Decode(payload.Content, &es))
	a.EqualValues(per.Etiquetas, es)

}

func TestQueryPersonaHandlers(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	addTestingHandlers(r)
	mock := test.NewMock("test", r)

	for _, fun := range []string{"QueryPersonaBasica", "QueryPersona"} {
		fun := fun
		t.Run(fun, func(t *testing.T) {
			a := assert.New(t)

			_, res, _, err := test.MockInvoke(t, mock, fun, 20104249729)
			a.NoError(err)
			a.EqualValues(200, res.Status)
			a.EqualValues("", res.Message)

			_, res, _, err = test.MockInvoke(t, mock, fun)
			a.NoError(err)
			a.EqualValues(400, res.Status)
			a.EqualValues("invalid argument: argument count mismatch: received 0 while expecting 1 (CUIT)", res.Message)

			_, res, _, err = test.MockInvoke(t, mock, fun, "-1")
			a.NoError(err)
			a.EqualValues(400, res.Status)
			a.EqualValues("invalid argument: CUIT argument 1: invalid natural integer: invalid syntax: '-1'", res.Message)

			pers := RandomPersonas(10, nil)
			for _, per := range pers {
				per := per
				t.Run(strconv.FormatUint(per.ID, 10), func(t *testing.T) {
					a := assert.New(t)

					_, res, _, err := test.MockInvoke(t, mock, "PutPersona", per)
					a.NoError(err)
					a.EqualValues(status.OK, res.Status)

					_, res, _, err = test.MockInvoke(t, mock, fun, per.ID)
					a.NoError(err)
					a.EqualValues(status.OK, res.Status)
				})
			}

		})
	}
}
