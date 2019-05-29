package personas_test

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/test"
)

func TestPutPersona(t *testing.T) {
	stub := test.SetInitTests(t)

	// Valid
	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}

	// Invalid cuit
	res = test.PutPersona(t, stub, 1)
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con un cuit invalido debe dar error")
		t.FailNow()
	}

	// distinct cuits
	var personaJSON = test.GetPersonaJSON(20255438795)
	res = stub.MockInvoke("1", [][]byte{[]byte("putPersona"), []byte("30679638943"), personaJSON})
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con cuits distintos debe dar error")
		t.FailNow()
	}
}

func TestPutPersonas(t *testing.T) {
	stub := test.SetInitTests(t)
	bs := []byte(`[` + string(test.GetPersonaJSON(20255438795)) + "," + string(test.GetPersonaJSON(30679638943)) + "]")
	res := stub.MockInvoke("1", [][]byte{[]byte("putPersonas"), bs})
	if res.Status != shim.OK {
		fmt.Println("putPersonas", res.Message)
		t.FailNow()
	} else {
		fmt.Println("putPersonas Ok!!!!")
	}

}
