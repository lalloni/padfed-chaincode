package inscripciones_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/test"
)

func TestPutPersonaImpuestos(t *testing.T) {
	stub := test.SetInitTests(t)

	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}

	impuestosJSON := `{"impuestos":[{"impuesto":30,"estado":"AC","periodo":199912},{"impuesto":32,"estado":"AC","periodo":199912}]}`

	res = stub.MockInvoke("1", [][]byte{[]byte("putPersonaImpuestos"), []byte("30679638943"), []byte(impuestosJSON)})
	if res.Status != shim.OK {
		fmt.Println("putPersonaImpuestos error", res.Message)
		t.FailNow()
	}
}
