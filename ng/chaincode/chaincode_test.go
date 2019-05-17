package chaincode_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
)

func TestInitHandler(t *testing.T) {
	a := assert.New(t)
	init := false
	mock := test.NewMock("cc", router.New(router.C(router.RH(func(*context.Context) *response.Response {
		init = true
		return response.OK("blah!")
	}))))
	res, p, err := test.MockInit(mock)
	a.NoError(err)
	a.True(init)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues("blah!", p.Content)
}

func TestInvokeHandler(t *testing.T) {
	a := assert.New(t)
	var (
		err  error
		call bool
		res  peer.Response
		p    response.Payload
	)
	r := router.New(nil)
	r.SetHandler(router.Name("f1"), authorization.Allowed, func(*context.Context) *response.Response {
		call = true
		return response.OK("blah!")
	})
	r.SetHandler(router.Name("f2"), authorization.Forbidden, func(*context.Context) *response.Response {
		call = true
		return response.OK("bleh!") // never invoked
	})

	mock := shim.NewMockStub("cc", chaincode.New("cc", r))

	call = false
	p = response.Payload{}
	res = mock.MockInvoke(uuid.New().String(), [][]byte{[]byte("f1")})
	a.True(call)
	a.EqualValues(http.StatusOK, res.Status)
	err = json.Unmarshal(res.Payload, &p)
	a.NoError(err)
	a.EqualValues("blah!", p.Content)

	call = false
	res = mock.MockInvoke(uuid.New().String(), [][]byte{[]byte("f2")})
	a.False(call)
	a.EqualValues(http.StatusForbidden, res.Status)
	a.NoError(err)
}
