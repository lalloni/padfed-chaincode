package personas_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/test"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestDelPersonaAssets(t *testing.T) {
	stub := test.SetInitTests(t)

	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		t.Fatalf("no se pudo cargar la persona: %v", res)
	}

	casos := []struct {
		s    string
		fail bool
	}{
		{`["PER_30679638943_IMP_30","PER_30679638943_IMP_124"]`, false},
		{`["PER_30679638943_IMP_30","PER_30679638943_IMP_30"]`, true},
		{`["PER_30679638943_IMP_30","PER_30679638945_IMP_124"]`, true},
		{`["PER_3067943_IMP_30","PER_30679638943_IMP_124"]`, true},
	}

	for _, caso := range casos {
		impuestosToDel := caso.s
		mustFail := caso.fail
		res = stub.MockInvoke("1", [][]byte{[]byte("delPersonaAssets"), []byte("30679638943"), []byte(impuestosToDel)})
		failed := res.Status != shim.OK
		if mustFail != failed {
			t.Errorf("failed: %v, expected: %v", failed, mustFail)
		}
	}

}
