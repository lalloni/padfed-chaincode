package personas_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/test"
)

func TestGetPersona(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)

	cuit := uint64(30679638943)

	personaPUT, err := canonicalize(test.GetPersonaJSON(cuit))
	a.NoError(err)

	res := stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaPUT})
	a.Equal(int32(shim.OK), res.Status, res.Payload)

	res = test.GetPersona(t, stub, cuit)
	a.Equal(int32(shim.OK), res.Status, res.Payload)

	personaGET, err := canonicalize(res.Payload)
	a.NoError(err)

	if !a.Equal(string(personaPUT), string(personaGET)) {
		t.Logf("expected:\n%s\nactual:\n%s", string(personaPUT), string(personaGET))
	}

}

func canonicalize(bs []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}
	a, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := json.Indent(b, a, "", "  "); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
