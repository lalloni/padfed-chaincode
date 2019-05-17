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

func TestGetPutStatesHandler(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)

	mock := test.NewMock("test", r.New(r.C(nil,
		r.R("gets", nil, generic.GetStatesHandler),
		r.R("puts", nil, generic.PutStatesHandler),
	)))

	res, payload, err := test.MockInvoke(mock, "puts", "key1", "foobarbaz")
	t.Logf("response status: %v message: %q payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(1, payload.Content)

	bs := make([]byte, 64)
	_, err = rand.Read(bs)
	a.NoError(err)
	res, payload, err = test.MockInvoke(mock, "puts", "key1", "foobarbaz", "key2", bs)
	t.Logf("response status: %v message: %q payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(2, payload.Content)

	res, payload, err = test.MockInvoke(mock, "gets", "key1")
	t.Logf("response status: %v message: %q payload: %s", res.Status, res.Message, string(res.Payload))
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues("foobarbaz", payload.Content)

	res, _, err = test.MockInvoke(mock, "gets", `["key1","key2"]`)
	t.Logf("response status: %v message: %q payload: %s", res.Status, res.Message, string(res.Payload))
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

}
