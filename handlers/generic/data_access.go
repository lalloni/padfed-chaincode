package generic

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func kvput(ctx *context.Context, key string, value []byte) *response.Response {
	err := ctx.Stub.PutState(key, value)
	if err != nil {
		return response.Error("putting state: %v", err)
	}
	return nil
}

func kvget(ctx *context.Context, key string) (*state, *response.Response) {
	bs, err := ctx.Stub.GetState(key)
	if err != nil {
		return nil, response.Error("getting key: %v", err)
	}
	return newstate(key, bs), nil
}

// TODO utilizar esta para implementar lectura de ranges
//nolint:deadcode,unused
func krget(ctx *context.Context, key1, key2 string) ([]*state, *response.Response) {
	it, err := ctx.Stub.GetStateByRange(key1, key2)
	if err != nil {
		return nil, response.Error("getting key range: %v", err)
	}
	defer it.Close()
	ss := []*state{}
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, response.Error("getting next key in range: %v", err)
		}
		ss = append(ss, newstate(kv.Key, kv.Value))
	}
	return ss, nil
}

func khget(ctx *context.Context, key string) ([]*statemod, *response.Response) {
	hi, err := ctx.Stub.GetHistoryForKey(key)
	if err != nil {
		return nil, response.Error("getting key history: %v", err)
	}
	defer hi.Close()
	mods := []*statemod{}
	for hi.HasNext() {
		km, err := hi.Next()
		if err != nil {
			return nil, response.Error("getting key modification: %v", err)
		}
		mods = append(mods, newstatemod(km))
	}
	return mods, nil
}

func rangekeys(ctx *context.Context, key1, key2 string) ([]string, *response.Response) {
	it, err := ctx.Stub.GetStateByRange(key1, key2)
	if err != nil {
		return nil, response.Error("getting key range: %v", err)
	}
	defer it.Close()
	ss := []string{}
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, response.Error("getting next key in range: %v", err)
		}
		ss = append(ss, kv.Key)
	}
	return ss, nil
}
