package impuesto

import (
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/handlers/common"
	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/impuesto"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(
		common.Defaults, // i.e. get, getrange, has, put, putlist, del, delrange
		common.WithIDParam(CodigoImpuestoParam),
		common.WithItemParam(ImpuestoParam),
		common.WithListParam(ImpuestoListParam),
	)
	if !testing {
		opts = append(opts, common.WithWriteCheck(common.AFIP))
	}
	common.AddCRUDHandlers(r, model.Schema, opts...)
}
