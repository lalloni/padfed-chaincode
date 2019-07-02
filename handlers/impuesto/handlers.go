package impuesto

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/common"
	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/impuesto"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/router"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(common.Defaults,
		common.WithIDParam(CodigoImpuestoParam),
		common.WithItemParam(ImpuestoParam),
		common.WithListParam(ImpuestoListParam),
	)
	if !testing {
		opts = append(opts, common.WithWriteCheck(common.AFIP))
	}
	common.AddCRUDHandlers(r, model.Schema, opts...)
}
