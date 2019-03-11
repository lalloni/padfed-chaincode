package inscripciones

import "strconv"

func GetImpuestoKeyByCuitID(cuit uint64, impuesto int32) string {
	cuitStr := strconv.FormatUint(cuit, 10)
	impStr := strconv.Itoa(int(impuesto))
	return "PER_" + cuitStr + "_IMP_" + impStr
}
