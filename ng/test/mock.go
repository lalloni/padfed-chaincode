package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func NewMock(name string, r router.Router) *shim.MockStub {
	return shim.NewMockStub(name, chaincode.New(name, "test", r))
}

func MockInvoke(t *testing.T, stub *shim.MockStub, function string, args ...interface{}) (string, *peer.Response, *response.Payload, error) {
	tx := uuid.New().String()
	aa := append([]interface{}{function}, args...)
	bs, err := arguments(aa)
	if err != nil {
		return "", nil, nil, err
	}
	f, _, err := context.ParseFunction([]byte(function))
	if err != nil {
		return "", nil, nil, err
	}
	res, payload, err := result(stub.MockInvoke(tx, bs))
	t.Logf("\n→ call function %q arguments: %s\n← response status: %v message: %q payload: %v error: %v", f, format(bs[1:]), res.Status, res.Message, strings.Trim(string(res.Payload), "\n"), err)
	return tx, res, payload, err
}

func MockInit(stub *shim.MockStub, args ...interface{}) (*peer.Response, *response.Payload, error) {
	bs, err := arguments(args)
	if err != nil {
		return nil, nil, err
	}
	return result(stub.MockInit(uuid.New().String(), bs))
}

func result(r peer.Response) (*peer.Response, *response.Payload, error) {
	p := response.Payload{}
	err := json.Unmarshal(r.Payload, &p)
	if err != nil {
		return nil, nil, err
	}
	return &r, &p, nil
}

func arguments(args []interface{}) ([][]byte, error) {
	bs := [][]byte{}
	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			bs = append(bs, v)
		case string:
			bs = append(bs, []byte(v))
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bs = append(bs, b)
		}
	}
	return bs, nil
}

func format(bss [][]byte) string {
	ss := []string{}
	for _, bs := range bss {
		ss = append(ss, fmt.Sprintf("%#v", string(bs)))
	}
	return strings.Join(ss, ",")
}
