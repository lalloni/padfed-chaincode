package impuestos

import (
	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/storeutil"
)

var Schema = store.MustPrepare(store.Composite{
	Name:            "impuesto",
	Creator:         func() interface{} { return &Impuesto{} },
	IdentifierField: "Codigo",
	IdentifierKey:   storeutil.Uint64Key("imp"),
	KeyIdentifier:   storeutil.Uint64Identifier(0),
	KeyBaseName:     "imp",
})
