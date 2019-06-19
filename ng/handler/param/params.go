package param

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var uint64t = reflect.TypeOf(uint64(0))

var Uint64 = New("natural integer", uint64t, parseuint64)

func parseuint64(args [][]byte) (interface{}, int, error) {
	s := string(args[0])
	r, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		if e, ok := err.(*strconv.NumError); ok {
			err = e.Err
		}
		return nil, 0, errors.Errorf("invalid natural integer: %v: '%v'", err, s)
	}
	return r, 1, nil
}

var stringt = reflect.TypeOf("")

var String = New("string", stringt, parsestring)

func parsestring(args [][]byte) (interface{}, int, error) {
	bs := args[0]
	if !utf8.Valid(bs) {
		return nil, 0, errors.Errorf("invalid UTF-8 string")
	}
	return string(bs), 1, nil
}
