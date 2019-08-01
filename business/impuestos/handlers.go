package impuestos

import (
	"github.com/lalloni/fabrikit/chaincode/handlerutil/crud"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(
		crud.Defaults, // i.e. get, getrange, has, put, putlist, del, delrange
		crud.WithIDParam(CodigoImpuestoParam),
		crud.WithItemParam(ImpuestoParam),
		crud.WithListParam(ImpuestoListParam),
	)
	if !testing {
		opts = append(opts, crud.WithWriteCheck(common.AFIP))
	}
	crud.AddHandlers(r, Schema, opts...)
}
