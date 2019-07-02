package impuesto

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var Schema = meta.MustPrepare(meta.Composite{
	Name:            "impuesto",
	Creator:         func() interface{} { return &Impuesto{} },
	IdentifierField: "Codigo",
	IdentifierKey:   common.Uint64Key("imp"),
	KeyIdentifier:   common.Uint64Identifier(0),
	KeyBaseName:     "imp",
})
