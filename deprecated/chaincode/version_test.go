package chaincode_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/chaincode"
)

func TestVersionHandler(t *testing.T) {
	version := "99.66.77.69"
	a := assert.New(t)
	mock := shim.NewMockStub("cc", chaincode.New(shim.NewLogger("cc"), version, true))
	res := mock.MockInvoke(uuid.New().String(), [][]byte{[]byte("version")})
	a.EqualValues(shim.OK, res.Status)
	a.EqualValues(version, res.Payload)
}

func TestErrorResponseVersion(t *testing.T) {
	version := "99.66.77.69"
	a := assert.New(t)
	mock := shim.NewMockStub("cc", chaincode.New(shim.NewLogger("cc"), version, true))
	res := mock.MockInvoke(uuid.New().String(), [][]byte{[]byte("putPersona"), []byte("")})
	a.EqualValues(shim.ERROR, res.Status)
	t.Logf("response was: %#v", res)
	fail := map[string]interface{}{}
	a.NoError(json.Unmarshal([]byte(res.Message), &fail)) // WTF?
	a.Equal(version, fail["version"])
}
