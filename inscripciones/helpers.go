package inscripciones

import (
	"strconv"
)

func GetImpuestoKeyByCuitID(cuit uint64, impuesto uint) string {
	return "PER_" + strconv.FormatUint(cuit, 10) + "_IMP_" + strconv.FormatUint(uint64(impuesto), 10)
}
