package chaincode_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/authorization"
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
	_, res, p, err := test.MockInit(t, mock)
	a.NoError(err)
	a.True(init)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues("blah!", p.Content)
}

func TestInvokeAuthorizationHandler(t *testing.T) {
	a := assert.New(t)

	var call bool
	mock := test.NewMock("cc", router.New(router.C(nil,
		router.R("f1", authorization.Allowed, func(*context.Context) *response.Response {
			call = true
			return response.OK("blah!")
		}),
		router.R("f2", authorization.Forbidden, func(*context.Context) *response.Response {
			call = true
			return response.OK("bleh!") // never invoked
		}),
	)))

	call = false
	_, res, p, err := test.MockInvoke(t, mock, "f1")
	a.NoError(err)
	a.True(call)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues("blah!", p.Content)

	call = false
	_, res, _, err = test.MockInvoke(t, mock, "f2")
	a.NoError(err)
	a.False(call)
	a.EqualValues(http.StatusForbidden, res.Status)
}

func TestDebugCall(t *testing.T) {
	a := assert.New(t)

	mock := test.NewMock("cc", router.New(router.C(nil,
		router.RNH("success", func(*context.Context) *response.Response {
			return response.OK("blah!")
		}),
		router.RNH("fail", func(*context.Context) *response.Response {
			return response.Error("bleh!")
		}),
	)))

	_, res, p, err := test.MockInvoke(t, mock, "success")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(&response.Payload{Content: "blah!"}, p)
	a.EqualValues("", res.Message)

	_, res, _, err = test.MockInvoke(t, mock, "fail")
	a.NoError(err)
	a.Nil(res.Payload)
	a.EqualValues(http.StatusInternalServerError, res.Status)
	a.EqualValues("bleh!", res.Message)

	tx, res, p, err := test.MockInvoke(t, mock, "success?debug")
	a.NoError(err)
	a.EqualValues(http.StatusOK, res.Status)
	a.EqualValues(&response.Payload{
		Content:     "blah!",
		Chaincode:   &response.Chaincode{Version: "test"},
		Transaction: &response.Transaction{ID: tx, Function: "success"},
	}, p)
	a.EqualValues("", res.Message)

	tx, res, p, err = test.MockInvoke(t, mock, "fail?debug")
	a.NoError(err)
	a.EqualValues(http.StatusInternalServerError, res.Status)
	a.EqualValues(&response.Payload{
		Chaincode:   &response.Chaincode{Version: "test"},
		Transaction: &response.Transaction{ID: tx, Function: "fail"},
	}, p)
	a.EqualValues("bleh!", res.Message)

}
