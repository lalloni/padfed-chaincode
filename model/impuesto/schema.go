package impuesto

import (
	"github.com/lalloni/fabrikit/chaincode/store"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/common"
)

var Schema = store.MustPrepare(store.Composite{
	Name:            "impuesto",
	Creator:         func() interface{} { return &Impuesto{} },
	IdentifierField: "Codigo",
	IdentifierKey:   common.Uint64Key("imp"),
	KeyIdentifier:   common.Uint64Identifier(0),
	KeyBaseName:     "imp",
})
