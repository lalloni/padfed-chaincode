package fabric_test

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestQueryByKey(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)
	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		t.Errorf("putPersona failed with: %s", res.Message)
	}
	res = test.QueryByKey(t, stub, key.Based("per", "30679638943").Tagged("per").String())
	if res.Status != shim.OK {
		t.Errorf("queryByKey failed with: %s", res.Message)
	}
	t.Logf("queryByKey result payload: %s", string(res.Payload))
	a.NotEqual("[]", string(res.Payload))
}
