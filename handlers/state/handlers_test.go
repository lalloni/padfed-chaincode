package state

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response/status"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestGetPutDelStatesHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("gets", nil, GetStatesHandler)
	r.SetHandler("dels", nil, DelStatesHandler)
	r.SetHandler("puts", nil, PutStatesHandler)

	mock := test.NewMock("test", r)

	bs := make([]byte, 64)
	_, err := rand.Read(bs)
	a.NoError(err)
	bs64 := base64.StdEncoding.EncodeToString(bs)

	t.Run("single put", func(t *testing.T) {
		a := assert.New(t)
		_, res, payload, err := test.MockInvoke(t, mock, "puts", "key1", "foobarbaz")
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues(1, payload.Content)
	})

	t.Run("double put", func(t *testing.T) {
		a := assert.New(t)
		_, res, payload, err := test.MockInvoke(t, mock, "puts", "key1", "foobarbaz", "key2", bs)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues(2, payload.Content)
	})

	t.Run("single point query", func(t *testing.T) {
		a := assert.New(t)
		_, res, payload, err := test.MockInvoke(t, mock, "gets", "key1")
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues("foobarbaz", payload.Content)
	})

	t.Run("multiple point query", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", `["key1","key2"]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 2)
		c0 := cs[0].Map()
		a.EqualValues("key1", c0["key"].Str)
		a.EqualValues("foobarbaz", c0["content"].Str)
		c1 := cs[1].Map()
		a.EqualValues("key2", c1["key"].Str)
		a.EqualValues("base64", c1["encoding"].Str)
		a.EqualValues(bs64, c1["content"].Str)
	})

	//nolint:dupl
	t.Run("single prefix range query", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", `[["key"]]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 1)
		c0 := cs[0].Array()
		a.Len(c0, 2)
		c00 := c0[0].Map()
		a.EqualValues("key1", c00["key"].Str)
		a.EqualValues("foobarbaz", c00["content"].Str)
		c01 := c0[1].Map()
		a.EqualValues("key2", c01["key"].Str)
		a.EqualValues("base64", c01["encoding"].Str)
		a.EqualValues(bs64, c01["content"].Str)
	})

	//nolint:dupl
	t.Run("single range query", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", `[["key0","key3"]]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 1)
		c0 := cs[0].Array()
		a.Len(c0, 2)
		c00 := c0[0].Map()
		a.EqualValues("key1", c00["key"].Str)
		a.EqualValues("foobarbaz", c00["content"].Str)
		c01 := c0[1].Map()
		a.EqualValues("key2", c01["key"].Str)
		a.EqualValues("base64", c01["encoding"].Str)
		a.EqualValues(bs64, c01["content"].Str)
	})

	t.Run("single left open range query", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", `[["","key2"]]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 1)
		c0 := cs[0].Array()
		a.Len(c0, 1)
		c00 := c0[0].Map()
		a.EqualValues("key1", c00["key"].Str)
		a.EqualValues("foobarbaz", c00["content"].Str)
	})

	t.Run("single right open range query", func(t *testing.T) {
		a := assert.New(t)
		t.Skip("mock stub no responde como se espera (volver a probar con fabric 1.4.1)")
		_, res, _, err := test.MockInvoke(t, mock, "gets", `[["key1",""]]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 1)
		c0 := cs[0].Array()
		a.Len(c0, 2)
		c00 := c0[0].Map()
		a.EqualValues("key1", c00["key"].Str)
		a.EqualValues("foobarbaz", c00["content"].Str)
		c01 := c0[1].Map()
		a.EqualValues("key2", c01["key"].Str)
		a.EqualValues("base64", c01["encoding"].Str)
		a.EqualValues(bs64, c01["content"].Str)
	})

	t.Run("multiple mixed queries", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", `[["key1","key2"],"key1",["key"],"key3"]`)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)

		cs := gjson.GetBytes(res.Payload, "content").Array()
		a.Len(cs, 4)

		c0 := cs[0].Array()
		a.Len(c0, 1)
		c00 := c0[0].Map()
		a.EqualValues("key1", c00["key"].Str)
		a.EqualValues("foobarbaz", c00["content"].Str)

		c1 := cs[1].Map()
		a.EqualValues("key1", c1["key"].Str)
		a.EqualValues("foobarbaz", c1["content"].Str)

		c2 := cs[2].Array()
		a.Len(c2, 2)
		c20 := c2[0].Map()
		a.EqualValues("key1", c20["key"].Str)
		a.EqualValues("foobarbaz", c20["content"].Str)
		c21 := c2[1].Map()
		a.EqualValues("key2", c21["key"].Str)
		a.EqualValues("base64", c21["encoding"].Str)
		a.EqualValues(bs64, c21["content"].Str)

		c3 := cs[3].Map()
		a.EqualValues("key3", c3["key"].Str)
		a.EqualValues(nil, c3["content"].Value())
	})

	t.Run("del key1", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "dels", "key1")
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	})

	t.Run("get key1 missing", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", "key1")
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)
	})

	t.Run("get key2", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", "key2")
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	})

	t.Run("del key2", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "dels", "key2")
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
	})

	t.Run("get key2 missing", func(t *testing.T) {
		a := assert.New(t)
		_, res, _, err := test.MockInvoke(t, mock, "gets", "key2")
		a.NoError(err)
		a.EqualValues(status.NotFound, res.Status)
	})

}

func TestGetStatesHistoryHandler(t *testing.T) {

	t.Skip("history not implemented in mock stub")

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()
	r.SetHandler("geth", nil, GetStatesHistoryHandler)
	r.SetHandler("puts", nil, PutStatesHandler)

	mock := test.NewMock("test", r)

	puts := func(key string, val string) {
		_, res, payload, err := test.MockInvoke(t, mock, "puts", key, val)
		a.NoError(err)
		a.EqualValues(status.OK, res.Status)
		a.EqualValues(1, payload.Content)
	}

	puts("key1", "foo")
	puts("key1", "bar")
	puts("key1", "baz")

	// get state history not implemented in shim.MockStub so... :'(
	_, res, payload, err := test.MockInvoke(t, mock, "geth", "key1")
	a.NoError(err)
	a.Regexp("getting key history: not implemented", res.Message)
	a.EqualValues(status.Error, res.Status)
	a.EqualValues(nil, payload.Content)

}
