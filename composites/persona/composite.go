package personas

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var Persona = meta.MustPrepare(meta.Composite{
	Name:    "persona",
	Creator: func() interface{} { return model.NewPersona() },
	Keyer:   func(id interface{}) *key.Key { return key.Based("per", strconv.FormatUint(id.(uint64), 10)) },
	Singletons: []meta.Singleton{
		{Tag: "per", Field: "Persona"},
	},
	Collections: []meta.Collection{
		{Tag: "act", Field: "Actividades"},
		{Tag: "imp", Field: "Impuestos"},
		{Tag: "dom", Field: "Domicilios"},
		{Tag: "tel", Field: "Telefonos"},
		{Tag: "jur", Field: "Jurisdicciones"},
		{Tag: "ema", Field: "Emails"},
		{Tag: "arc", Field: "Archivos"},
		{Tag: "cat", Field: "Categorias"},
		{Tag: "eti", Field: "Etiquetas"},
		{Tag: "con", Field: "Contribuciones"},
		{Tag: "rel", Field: "Relaciones"},
	},
})
