package state

import (
	"github.com/pkg/errors"

	"github.com/hyperledger/fabric/protos/common"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response/status"
)

func setKeyValue(ctx *context.Context, key string, value []byte) error {
	err := ctx.Stub.PutState(key, value)
	if err != nil {
		return errors.Wrap(err, "putting state")
	}
	return nil
}

func keyState(ctx *context.Context, key string) (*state, error) {
	bs, err := ctx.Stub.GetState(key)
	if err != nil {
		return nil, errors.Wrap(err, "getting key")
	}
	return newstate(key, bs), nil
}

func keyRangeStates(ctx *context.Context, key1, key2 string) ([]*state, error) {
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

func keyHistory(ctx *context.Context, key string) ([]*statemod, error) {
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
		b, err := txBlock(ctx, km.GetTxId())
		if err != nil {
			return nil, err
		}
		mods = append(mods, newstatemod(km, b))
	}
	return mods, nil
}

var base [][]byte

func txBlock(ctx *context.Context, txid string) (uint64, error) {
	if len(base) == 0 {
		base = [][]byte{[]byte("GetBlockByTxID"), []byte(ctx.Stub.GetChannelID())}
	}
	res := ctx.Stub.InvokeChaincode("qscc", append(base, []byte(txid)), "")
	if res.GetStatus() != status.OK {
		return 0, errors.Errorf("getting transaction block number: %v", res.Message)
	}
	b := &common.Block{}
	err := b.XXX_Unmarshal(res.Payload)
	if err != nil {
		return 0, errors.Wrap(err, "parsing transaction block number")
	}
	return b.Header.Number, nil
}
