package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/chaincode"
)

func CheckInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", res.Message)
		t.FailNow()
	}
}

func CheckState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func PutPersona(t *testing.T, stub *shim.MockStub, cuit uint64) peer.Response {
	var personaJSON = GetPersonaJSON(cuit)
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("putPersona"), []byte(cuitStr), []byte(personaJSON)})
}

func QueryPersona(t *testing.T, stub *shim.MockStub, cuit uint64) peer.Response {
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("queryPersona"), []byte(cuitStr)})
}

func SetInitTests(t *testing.T) *shim.MockStub {
	scc := chaincode.New(shim.NewLogger("padfedcc"), true)
	stub := shim.NewMockStub("padfed", scc)
	CheckInit(t, stub, [][]byte{})
	return stub
}
