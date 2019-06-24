package meta

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var Persona = meta.MustPrepare(meta.Composite{
	Name:            "persona",
	KeyBaseName:     "per",
	IdentifierField: "ID",
	IdentifierKey:   common.Uint64Key("per"),
	KeyIdentifier:   common.Uint64Identifier(0),
	Creator:         func() interface{} { return &model.Persona{} },
	Singletons: []meta.Singleton{
		{Tag: "per", Field: "Persona"},
	},
	Collections: []meta.Collection{
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
