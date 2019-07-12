package persona

import (
	"github.com/lalloni/fabrikit/chaincode/store"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/common"
)

var Schema = store.MustPrepare(store.Composite{
	Name:            "persona",
	KeyBaseName:     "per",
	IdentifierField: "ID",
	IdentifierKey:   common.Uint64Key("per"),
	KeyIdentifier:   common.Uint64Identifier(0),
	Creator:         func() interface{} { return &Persona{} },
	Singletons: []store.Singleton{
		{Tag: "per", Field: "Persona"},
	},
	Collections: []store.Collection{
		{Tag: "act", Field: "Actividades"},
		{Tag: "imp", Field: "Impuestos"},
		{Tag: "dom", Field: "Domicilios"},
		{Tag: "dor", Field: "DomiciliosRoles"},
		{Tag: "tel", Field: "Telefonos"},
		{Tag: "jur", Field: "Jurisdicciones"},
		{Tag: "ema", Field: "Emails"},
		{Tag: "arc", Field: "Archivos"},
		{Tag: "cat", Field: "Categorias"},
		{Tag: "eti", Field: "Etiquetas"},
		{Tag: "con", Field: "Contribuciones"},
		{Tag: "rel", Field: "Relaciones"},
		{Tag: "cms", Field: "CMSedes"},
	},
})
