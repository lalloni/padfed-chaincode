package ranges

import (
	"reflect"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/handler/param"
)

var Param = param.New("key ranges", reflect.TypeOf(&Ranges{}), parser)

func parser(args [][]byte) (interface{}, int, error) {
	r, err := Parse(args[0])
	if err != nil {
		return nil, 0, err
	}
	return r, 1, nil
}
