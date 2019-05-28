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

	bs := make([]byte, 64)
	_, err := rand.Read(bs)
	a.NoError(err)

	t.Run("single put", func(t *testing.T) {
		res, payload, err := test.MockInvoke(t, mock, "puts", "key1", "foobarbaz")
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		a.EqualValues(1, payload.Content)
	})

	t.Run("double put", func(t *testing.T) {
		res, payload, err := test.MockInvoke(t, mock, "puts", "key1", "foobarbaz", "key2", bs)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		a.EqualValues(2, payload.Content)
	})

	t.Run("single point query", func(t *testing.T) {
		res, payload, err := test.MockInvoke(t, mock, "gets", "key1")
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		a.EqualValues("foobarbaz", payload.Content)
	})

	t.Run("multiple point query", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", `["key1","key2"]`)
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
	})

	//nolint:dupl
	t.Run("single prefix range query", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", `[["key"]]`)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		p, err := simplejson.NewJson(res.Payload)
		a.NoError(err)
		a.EqualValues("key1", p.Get("content").GetIndex(0).GetIndex(0).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).GetIndex(0).Get("content").MustString())
		a.EqualValues("key2", p.Get("content").GetIndex(0).GetIndex(1).Get("key").MustString())
		a.EqualValues("base64", p.Get("content").GetIndex(0).GetIndex(1).Get("encoding").MustString())
		bss, err := base64.StdEncoding.DecodeString(p.Get("content").GetIndex(0).GetIndex(1).Get("content").MustString())
		a.NoError(err)
		a.EqualValues(bs, bss)
	})

	//nolint:dupl
	t.Run("single range query", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", `[["key0","key3"]]`)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		p, err := simplejson.NewJson(res.Payload)
		a.NoError(err)
		a.EqualValues("key1", p.Get("content").GetIndex(0).GetIndex(0).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).GetIndex(0).Get("content").MustString())
		a.EqualValues("key2", p.Get("content").GetIndex(0).GetIndex(1).Get("key").MustString())
		a.EqualValues("base64", p.Get("content").GetIndex(0).GetIndex(1).Get("encoding").MustString())
		bss, err := base64.StdEncoding.DecodeString(p.Get("content").GetIndex(0).GetIndex(1).Get("content").MustString())
		a.NoError(err)
		a.EqualValues(bs, bss)
	})

	t.Run("single left open range query", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", `[["","key2"]]`)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		p, err := simplejson.NewJson(res.Payload)
		a.NoError(err)
		a.EqualValues("key1", p.Get("content").GetIndex(0).GetIndex(0).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).GetIndex(0).Get("content").MustString())
	})

	t.Run("single right open range query", func(t *testing.T) {
		// TODO caso comentado porque el mock no responde como se espera (revisar con fabric 1.4.1)
		//res, _, err = test.MockInvoke(t, mock, "gets", `[["key1",""]]`)
		//a.NoError(err)
		//a.EqualValues(http.StatusOK, res.Status)
		//p, err = simplejson.NewJson(res.Payload)
		//a.NoError(err)
		//a.EqualValues("key1", p.Get("content").GetIndex(0).GetIndex(0).Get("key").MustString())
		//a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).GetIndex(0).Get("content").MustString())
		//a.EqualValues("key2", p.Get("content").GetIndex(0).GetIndex(1).Get("key").MustString())
		//a.EqualValues("base64", p.Get("content").GetIndex(0).GetIndex(1).Get("encoding").MustString())
		//bss, err = base64.StdEncoding.DecodeString(p.Get("content").GetIndex(0).GetIndex(1).Get("content").MustString())
		//a.NoError(err)
		//a.EqualValues(bs, bss)
	})

	t.Run("multiple mixed queries", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", `[["key1","key2"],"key1",["key"],"key3"]`)
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
		p, err := simplejson.NewJson(res.Payload)
		a.NoError(err)
		// range query at 0
		a.Len(p.Get("content").GetIndex(0).MustArray(), 1)
		a.EqualValues("key1", p.Get("content").GetIndex(0).GetIndex(0).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(0).GetIndex(0).Get("content").MustString())
		// point query at 1
		a.EqualValues("key1", p.Get("content").GetIndex(1).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(1).Get("content").MustString())
		// prefix range query at 2
		a.Len(p.Get("content").GetIndex(2).MustArray(), 2)
		a.EqualValues("key1", p.Get("content").GetIndex(2).GetIndex(0).Get("key").MustString())
		a.EqualValues("foobarbaz", p.Get("content").GetIndex(2).GetIndex(0).Get("content").MustString())
		a.EqualValues("key2", p.Get("content").GetIndex(2).GetIndex(1).Get("key").MustString())
		a.EqualValues("base64", p.Get("content").GetIndex(2).GetIndex(1).Get("encoding").MustString())
		bss, err := base64.StdEncoding.DecodeString(p.Get("content").GetIndex(2).GetIndex(1).Get("content").MustString())
		a.NoError(err)
		a.EqualValues(bs, bss)
		// key3 does not exist at 3
		a.EqualValues("key3", p.Get("content").GetIndex(3).Get("key").MustString())
		_, present := p.Get("content").GetIndex(3).CheckGet("content")
		a.False(present)
	})

	t.Run("del key1", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "dels", "key1")
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	})

	t.Run("get key1 missing", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", "key1")
		a.NoError(err)
		a.EqualValues(http.StatusNotFound, res.Status)
	})

	t.Run("get key2", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", "key2")
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	})

	t.Run("del key2", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "dels", "key2")
		a.NoError(err)
		a.EqualValues(http.StatusOK, res.Status)
	})

	t.Run("get key2 missing", func(t *testing.T) {
		res, _, err := test.MockInvoke(t, mock, "gets", "key2")
		a.NoError(err)
		a.EqualValues(http.StatusNotFound, res.Status)
	})

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
