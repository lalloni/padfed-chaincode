package personas_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestQueryPersona(t *testing.T) {
	stub := test.SetInitTests(t)
	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}
	res = test.QueryPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("queryPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}
	fmt.Println("queryPersona ", string(res.Payload))
}
