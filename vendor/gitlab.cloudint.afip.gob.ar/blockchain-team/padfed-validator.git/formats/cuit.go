package formats

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/lalloni/afip/cuit"
)

var Cuit = FormatCheckerFunc(checkCuit)

var zero = big.NewRat(0, 1)

func checkCuit(input interface{}) bool {
	switch v := input.(type) {
	case *big.Rat:
		if !v.IsInt() {
			return false
		}
		if v.Cmp(zero) < 1 {
			return false
		}
		if v.Num().IsUint64() {
			return cuit.IsValid(v.Num().Uint64())
		}
		return false
	case uint:
		if strconv.IntSize < 64 {
			return false
		}
		return cuit.IsValid(uint64(v))
	case int:
		if strconv.IntSize < 64 {
			return false
		}
		if v < 0 {
			return false
		}
		return cuit.IsValid(uint64(v))
	case int64:
		if v < 0 {
			return false
		}
		return cuit.IsValid(uint64(v))
	case uint64:
		return cuit.IsValid(v)
	case string:
		i, err := cuit.Parse(v)
		return err != nil && cuit.IsValid(i)
	default:
		i, err := cuit.Parse(fmt.Sprint(input))
		return err != nil && cuit.IsValid(i)
	}
}
