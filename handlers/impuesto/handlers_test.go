package impuesto

import (
	"encoding/json"
	"sort"
	"testing"
	"unicode/utf8"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/response/status"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/impuesto"
)

func TestPutImpuesto(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)

	imp1 := &model.Impuesto{
		Codigo:      1,
		Org:         1,
		Abreviatura: "IVA",
		Nombre:      "Impuesto al valor agregado",
	}

	_, res, payload, err := test.MockInvoke(t, mock, "PutImpuesto", imp1)
	a.NoError(err)
	a.Empty(payload)
	a.EqualValues(status.OK, res.Status)

	imp2 := &model.Impuesto{}
	bs, err := mock.GetState("imp:1")
	a.NoError(err)
	err = json.Unmarshal(bs, imp2)
	a.NoError(err)
	a.EqualValues(imp1, imp2)

	bs, err = mock.GetState("imp:1#wit")
	a.NoError(err)
	a.EqualValues("1", string(bs))

	q, err := mock.GetStateByRange(prefix("imp:1"))
	a.NoError(err)
	ks := []string(nil)
	for q.HasNext() {
		n, err := q.Next()
		a.NoError(err)
		ks = append(ks, n.Key)
	}
	q.Close()
	sort.Strings(ks)
	a.EqualValues([]string{"imp:1", "imp:1#wit"}, ks)

}

func TestPutImpuestoList(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)

	imps := []*model.Impuesto{
		{
			Codigo:      1,
			Org:         1,
			Abreviatura: "IVA",
			Nombre:      "Impuesto al valor agregado",
		}, {
			Codigo:      2,
			Org:         1,
			Abreviatura: "GAN",
			Nombre:      "Impuesto a las ganancias",
		},
	}

	_, res, payload, err := test.MockInvoke(t, mock, "PutImpuestoList", imps)
	a.NoError(err)
	a.EqualValues(2, payload.Content)
	a.EqualValues(status.OK, res.Status)

	imp := &model.Impuesto{}
	bs, err := mock.GetState("imp:1")
	a.NoError(err)
	err = json.Unmarshal(bs, imp)
	a.NoError(err)
	a.EqualValues(imps[0], imp)

	bs, err = mock.GetState("imp:1#wit")
	a.NoError(err)
	a.EqualValues("1", string(bs))

	imp = &model.Impuesto{}
	bs, err = mock.GetState("imp:2")
	a.NoError(err)
	err = json.Unmarshal(bs, imp)
	a.NoError(err)
	a.EqualValues(imps[1], imp)

	bs, err = mock.GetState("imp:2#wit")
	a.NoError(err)
	a.EqualValues("1", string(bs))

	q, err := mock.GetStateByRange(prefix("imp:"))
	a.NoError(err)
	ks := []string(nil)
	for q.HasNext() {
		n, err := q.Next()
		a.NoError(err)
		ks = append(ks, n.Key)
	}
	q.Close()
	sort.Strings(ks)
	a.EqualValues([]string{"imp:1", "imp:1#wit", "imp:2", "imp:2#wit"}, ks)

}

func TestGetImpuesto(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)
	test.MockTransactionStart(t, mock)

	imp1 := &model.Impuesto{
		Codigo:      1,
		Org:         1,
		Abreviatura: "IVA",
		Nombre:      "Impuesto al valor agregado",
	}

	bs, err := json.Marshal(imp1)
	a.NoError(err)
	err = mock.PutState("imp:1", bs)
	a.NoError(err)
	err = mock.PutState("imp:1#wit", []byte("1"))
	a.NoError(err)

	_, res, payload, err := test.MockInvoke(t, mock, "GetImpuesto", 1)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	imp2 := &model.Impuesto{}
	err = mapstructure.Decode(payload.Content, imp2)
	a.NoError(err)
	a.EqualValues(imp1, imp2)

}

func TestHasImpuesto(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)

	test.MockTransactionStart(t, mock)

	imp1 := &model.Impuesto{
		Codigo:      1,
		Org:         1,
		Abreviatura: "IVA",
		Nombre:      "Impuesto al valor agregado",
	}

	bs, err := json.Marshal(imp1)
	a.NoError(err)
	err = mock.PutState("imp:1", bs)
	a.NoError(err)
	err = mock.PutState("imp:1#wit", []byte("1"))
	a.NoError(err)

	_, res, payload, err := test.MockInvoke(t, mock, "HasImpuesto", 1)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(true, payload.Content)

}

func TestDelImpuesto(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)

	test.MockTransactionStart(t, mock)

	imp1 := &model.Impuesto{
		Codigo:      1,
		Org:         1,
		Abreviatura: "IVA",
		Nombre:      "Impuesto al valor agregado",
	}

	bs, err := json.Marshal(imp1)
	a.NoError(err)
	err = mock.PutState("imp:1", bs)
	a.NoError(err)
	err = mock.PutState("imp:1#wit", []byte("1"))
	a.NoError(err)

	_, res, payload, err := test.MockInvoke(t, mock, "DelImpuesto", 1)
	a.NoError(err)
	a.Empty(payload)
	a.EqualValues(status.OK, res.Status)

	q, err := mock.GetStateByRange(prefix("imp:"))
	a.NoError(err)
	a.False(q.HasNext())
	q.Close()

}

func TestGetImpuestoAll(t *testing.T) {

	shim.SetLoggingLevel(shim.LogDebug)
	a := assert.New(t)

	r := router.New()

	addTestHandlers(r)

	mock := test.NewMock("test", r)

	imps := []*model.Impuesto{
		{
			Codigo:      1,
			Org:         1,
			Abreviatura: "IVA",
			Nombre:      "Impuesto al valor agregado",
		}, {
			Codigo:      2,
			Org:         1,
			Abreviatura: "GAN",
			Nombre:      "Impuesto a las ganancias",
		},
	}

	_, res, payload, err := test.MockInvoke(t, mock, "PutImpuestoList", imps)
	a.NoError(err)
	a.EqualValues(2, payload.Content)
	a.EqualValues(status.OK, res.Status)

	_, res, payload, err = test.MockInvoke(t, mock, "GetImpuestoAll")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.NotEmpty(payload)
	imps2 := []*model.Impuesto{}
	err = mapstructure.Decode(payload.Content, &imps2)
	a.NoError(err)
	a.EqualValues(imps, imps2)

}

func prefix(s string) (string, string) {
	return s, s + string(utf8.MaxRune)
}
