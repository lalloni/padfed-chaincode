package personas_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestDelPersonasByRange(t *testing.T) {
	stub := test.SetInitTests(t)

	cuils := []uint64{20066042333, 20066675573, 20066806163, 20068854785, 20176058650}
	// creo varias personas
	for _, cuil := range cuils {
		res := test.PutPersona(t, stub, cuil)
		if res.Status != shim.OK {
			fmt.Println("putPersona", "cuit", "failed", res.Message)
			t.FailNow()
		}
	}

	cuilsRestantes := 2
	res := stub.MockInvoke("1", [][]byte{[]byte("delPersonasByRange"), []byte("20066600000"), []byte("20068900000")})
	if res.Status != shim.OK {
		fmt.Println("delPersonasByRange", "cuit", "failed", res.Message)
		t.FailNow()
	}

	fmt.Println("--- ESTADO FINAL ---")
	res = stub.MockInvoke("1", [][]byte{[]byte("queryAllPersona")})
	if res.Status != shim.OK {
		fmt.Println("queryAllPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}
	v := []interface{}{}
	err := json.Unmarshal(res.Payload, &v)
	if err != nil {
		t.Errorf("unmarshalling response message: %v", err)
	}
	t.Logf("unmarshalled: %+v", v)
	if len(v) != cuilsRestantes {
		t.Errorf("no coincide")
	}
}
