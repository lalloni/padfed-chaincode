package test

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/chaincode"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func NewMock(name string, r router.Router) *shim.MockStub {
	return shim.NewMockStub(name, chaincode.New(name, r))
}

func MockInvoke(stub *shim.MockStub, function string, args ...interface{}) (*peer.Response, *response.Payload, error) {
	aa := append([]interface{}{function}, args...)
	bs, err := arguments(aa)
	if err != nil {
		return nil, nil, err
	}
	return result(stub.MockInvoke(uuid.New().String(), bs))
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
