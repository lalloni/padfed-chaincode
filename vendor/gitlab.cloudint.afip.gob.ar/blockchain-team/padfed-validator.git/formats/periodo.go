package formats

import (
	"fmt"
	"math/big"
	"math/bits"

	"github.com/lalloni/afip/periodo"
)

var PeriodoDiario = FormatCheckerFunc(checkPeriodoDiario)
var PeriodoMensual = FormatCheckerFunc(checkPeriodoMensual)
var PeriodoAnual = FormatCheckerFunc(checkPeriodoAnual)

var z = big.NewRat(0, 1)

func checkPeriodoDiario(in interface{}) bool {
	switch v := in.(type) {
	case *big.Rat:
		u, ok := convert(v)
		if !ok {
			return false
		}
		return periodo.CheckPeriodoDiarioCompound(u)
	case string:
		ok, _, _, _ := periodo.Parse(periodo.Diario, v)
		return ok
	default:
		ok, _, _, _ := periodo.Parse(periodo.Mensual, fmt.Sprint(in))
		return ok
	}
}

func checkPeriodoMensual(in interface{}) bool {
	switch v := in.(type) {
	case *big.Rat:
		u, ok := convert(v)
		if !ok {
			return false
		}
		return periodo.CheckPeriodoMensualCompound(u)
	case string:
		ok, _, _, _ := periodo.Parse(periodo.Mensual, v)
		return ok
	default:
		ok, _, _, _ := periodo.Parse(periodo.Mensual, fmt.Sprint(in))
		return ok
	}
}

func checkPeriodoAnual(in interface{}) bool {
	switch v := in.(type) {
	case *big.Rat:
		u, ok := convert(v)
		if !ok {
			return false
		}
		return periodo.CheckPeriodoAnual(u)
	case string:
		ok, _, _, _ := periodo.Parse(periodo.Anual, v)
		return ok
	default:
		ok, _, _, _ := periodo.Parse(periodo.Anual, fmt.Sprint(in))
		return ok
	}
}

func convert(r *big.Rat) (uint, bool) {
	if r == nil || !r.IsInt() || r.Cmp(z) == -1 || !r.Num().IsUint64() {
		return 0, false
	}
	u := r.Num().Uint64()
	if bits.Len64(u) > bits.UintSize {
		return 0, false
	}
	return uint(u), true
}
