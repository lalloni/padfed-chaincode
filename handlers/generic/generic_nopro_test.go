package generic_test

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/generic"
	r "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestGetPutDelStatesHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("gets", nil, generic.GetStatesHandler),
		r.R("puts", nil, generic.PutStatesHandler),
		r.R("dels", nil, generic.DelStatesHandler),
	)))

	res, payload, err := test.MockInvoke(t, mock, "puts", "key1", "foobarbaz")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(1, payload.Content)

	bs := make([]byte, 64)
	_, err = rand.Read(bs)
	a.NoError(err)
	res, payload, err = test.MockInvoke(t, mock, "puts", "key1", "foobarbaz", "key2", bs)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(2, payload.Content)

	res, payload, err = test.MockInvoke(t, mock, "gets", "key1")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues("foobarbaz", payload.Content)

	res, _, err = test.MockInvoke(t, mock, "gets", `["key1","key2"]`)
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	p, err := simplejson.NewJson(res.Payload)
	a.NoError(err)
	a.EqualValues("key1", p.Get("content").GetIndex(0).Get("key").MustString())
	a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).Get("content").MustString())
	a.EqualValues("key2", p.Get("content").GetIndex(1).Get("key").MustString())
	a.EqualValues("base64", p.Get("content").GetIndex(1).Get("encoding").MustString())
	bss, err := base64.StdEncoding.DecodeString(p.Get("content").GetIndex(1).Get("content").MustString())
	a.NoError(err)
	a.EqualValues(bs, bss)

	res, _, err = test.MockInvoke(t, mock, "dels", "key1")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)

	res, _, err = test.MockInvoke(t, mock, "gets", "key1")
	a.NoError(err)
	a.EqualValues(http.StatusNotFound, res.Status)

	res, _, err = test.MockInvoke(t, mock, "gets", "key2")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)

	res, _, err = test.MockInvoke(t, mock, "dels", "key2")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)

	res, _, err = test.MockInvoke(t, mock, "gets", "key2")
	a.NoError(err)
	a.EqualValues(http.StatusNotFound, res.Status)

}

func TestGetStatesHistoryHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("geth", nil, generic.GetStatesHistoryHandler),
		r.R("puts", nil, generic.PutStatesHandler),
	)))

	puts := func(key string, val string) {
		res, payload, err := test.MockInvoke(t, mock, "puts", key, val)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		a.EqualValues(1, payload.Content)
	}

	puts("key1", "foo")
	puts("key1", "bar")
	puts("key1", "baz")

	// get state history not implemented in shim.MockStub so... :'(
	res, payload, err := test.MockInvoke(t, mock, "geth", "key1")
	a.NoError(err)
	a.Regexp("getting key history: not implemented", res.Message)
	a.EqualValues(http.StatusInternalServerError, res.Status)
	a.EqualValues(nil, payload.Content)

}
