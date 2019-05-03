package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/chaincode"
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

func QueryByKey(t *testing.T, stub *shim.MockStub, key string) peer.Response {
	return stub.MockInvoke("1", [][]byte{[]byte("queryByKey"), []byte(key)})
}

func PutPersona(t *testing.T, stub *shim.MockStub, cuit uint64) peer.Response {
	var personaJSON = GetPersonaJSON(cuit)
	return stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaJSON})
}

func GetPersona(t *testing.T, stub *shim.MockStub, cuit uint64) peer.Response {
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("getPersona"), []byte(cuitStr)})
}

func QueryPersona(t *testing.T, stub *shim.MockStub, cuit uint64) peer.Response {
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("queryPersona"), []byte(cuitStr)})
}

func SetInitTests(t *testing.T) *shim.MockStub {
	scc := chaincode.New(shim.NewLogger("padfedcc"), "test", true)
	stub := shim.NewMockStub("padfed", scc)
	CheckInit(t, stub, [][]byte{})
	return stub
}

func GetPersonaJSON(cuit uint64) []byte {
	if cuit >= 30000000000 {
		return GetPersonaJurídica(cuit)
	}
	return GetPersonaFísica(cuit)
}
