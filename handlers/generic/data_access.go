package generic

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

func kvput(ctx *context.Context, key string, value []byte) error {
	err := ctx.Stub.PutState(key, value)
	if err != nil {
		return errors.Wrap(err, "putting state")
	}
	return nil
}

func kvget(ctx *context.Context, key string) (*state, error) {
	bs, err := ctx.Stub.GetState(key)
	if err != nil {
		return nil, errors.Wrap(err, "getting key")
	}
	return newstate(key, bs), nil
}

func krget(ctx *context.Context, key1, key2 string) ([]*state, error) {
	it, err := ctx.Stub.GetStateByRange(key1, key2)
	if err != nil {
		return nil, errors.Wrap(err, "getting key range")
	}
	defer it.Close()
	ss := []*state{}
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, errors.Wrap(err, "getting next key in range")
		}
		ss = append(ss, newstate(kv.Key, kv.Value))
	}
	return ss, nil
}

func khget(ctx *context.Context, key string) ([]*statemod, error) {
	hi, err := ctx.Stub.GetHistoryForKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "getting key history")
	}
	defer hi.Close()
	mods := []*statemod{}
	for hi.HasNext() {
		km, err := hi.Next()
		if err != nil {
			return nil, errors.Wrap(err, "getting key modification")
		}
		mods = append(mods, newstatemod(km))
	}
	return mods, nil
}
