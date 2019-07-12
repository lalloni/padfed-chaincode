package state

import (
	"reflect"

	"github.com/lalloni/fabrikit/chaincode/handler/param"
)

var RangesParam = param.Typed("key ranges", reflect.TypeOf(&Ranges{}), parser)

func parser(arg []byte) (interface{}, error) {
	r, err := Parse(arg)
	if err != nil {
		return nil, err
	}
	return r, nil
}
