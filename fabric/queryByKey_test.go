package fabric_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestQueryByKey(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)
	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}
	res = test.QueryByKey(t, stub, "PER_30679638943")
	if res.Status != shim.OK {
		fmt.Println("queryByKey", "failed", res.Message)
		t.FailNow()
	}
	fmt.Println("queryByKey ", string(res.Payload))
	a.NotEqual(string(res.Payload), "[]", "resultado no esperado")
}
