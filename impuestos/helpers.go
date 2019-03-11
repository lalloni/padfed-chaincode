package impuestos

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func HasDuplicatedImpuestos(impuestos []*model.Impuesto) (bool, *model.Impuesto) {
	var index = make(map[int]*model.Impuesto)
	for _, imp := range impuestos {
		if _, exist := index[int(imp.Impuesto)]; exist {
			return exist, imp
		}
		index[int(imp.Impuesto)] = imp
	}
	return false, &model.Impuesto{}
}
