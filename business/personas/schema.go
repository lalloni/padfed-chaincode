package personas

import (
	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/storeutil"
)

var Schema = store.MustPrepare(store.Composite{
	Name:            "persona",
	KeyBaseName:     "per",
	IdentifierField: "ID",
	IdentifierKey:   storeutil.Uint64Key("per"),
	KeyIdentifier:   storeutil.Uint64Identifier(0),
	Creator:         func() interface{} { return &Persona{} },
	Singletons: []store.Singleton{
		{Tag: "per", Name: "Basica", Field: "Persona"},
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
