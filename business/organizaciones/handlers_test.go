package organizaciones

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/test"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestGetAllHandler(t *testing.T) {

	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	AddHandlers(r)

	mock := test.NewMock("test", r)

	_, _, payload, err := test.MockInvoke(t, mock, "GetOrganizacionAll")
	a.NoError(err)

	orgs := []*Org(nil)
	err = mapstructure.Decode(payload.Content, &orgs)
	a.NoError(err)

	a.EqualValues(GetAll(), orgs)
}

func TestGetHandler(t *testing.T) {

	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)

	r := router.New()

	AddHandlers(r)

	mock := test.NewMock("test", r)

	_, _, payload, err := test.MockInvoke(t, mock, "GetOrganizacion", 1)
	a.NoError(err)

	org := &Org{}
	err = mapstructure.Decode(payload.Content, &org)
	a.NoError(err)

	a.EqualValues(GetByID(1), org)
	a.EqualValues(1, org.ID)
}
