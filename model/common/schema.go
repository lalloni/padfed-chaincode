package common

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

func Uint64Key(name string) meta.KeyFunc {
	return func(id interface{}) (*key.Key, error) {
		return key.NewBase(name, strconv.FormatUint(id.(uint64), 10)), nil
	}
}

func Uint64Identifier(seg int) meta.ValFunc {
	return func(k *key.Key) (interface{}, error) {
		return strconv.ParseUint(k.Base[seg].Value, 10, 64)
	}
}
