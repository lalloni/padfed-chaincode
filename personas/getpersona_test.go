package personas_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/bitly/go-simplejson"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
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

func TestGetPersonaLenientErrors(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)

	cuit := uint64(30679638943)

	personaPUT, err := canonicalize(test.GetPersonaJSON(cuit))
	a.NoError(err)

	res := stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaPUT})
	a.EqualValues(shim.OK, res.Status)

	// modifico y guardo la personabasica inválida
	json, err := simplejson.NewJson(personaPUT)
	a.NoError(err)
	json = json.Get("persona")
	json.Set("pais", "un no número")
	bs, err := json.MarshalJSON()
	a.NoError(err)
	stub.MockTransactionStart("1")
	err = stub.PutState("per:30679638943#per", bs)
	a.NoError(err)
	stub.MockTransactionEnd("1")

	// cargo en modo lenient pidiendo que embeba errores
	res = stub.MockInvoke("1", [][]byte{[]byte("getPersona?lenientread&embederrors"), []byte(strconv.FormatUint(cuit, 10))})
	a.EqualValues(shim.OK, res.Status)

	// controlo que vengan errores
	json, err = simplejson.NewJson(res.Payload)
	a.NoError(err)
	ee, err := json.Get("errors").Array()
	a.NoError(err)
	a.NotEmpty(ee)

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
