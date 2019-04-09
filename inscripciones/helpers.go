package inscripciones

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func GetImpuestoKeyByCuitID(cuit uint64, impuesto uint) string {
	return "PER_" + strconv.FormatUint(cuit, 10) + "_IMP_" + strconv.FormatUint(uint64(impuesto), 10)
}

func HasDuplicatedImpuestos(impuestos []*model.PersonaImpuesto) (bool, *model.PersonaImpuesto) {
	index := map[uint]struct{}{}
	for _, imp := range impuestos {
		if _, ok := index[imp.Impuesto]; ok {
			return true, imp
		}
		index[imp.Impuesto] = struct{}{}
	}
	return false, nil
}
