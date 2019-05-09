package cast

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var (
	Persona = meta.MustPrepare(
		meta.Composite{
			Name:            "persona",
			IdentifierField: "ID",
			Creator: func() interface{} {
				return model.NewPersona()
			},
			IdentifierKey: func(id interface{}) *key.Key {
				return key.NewBase("per", strconv.FormatUint(id.(uint64), 10))
			},
			KeyIdentifier: func(k *key.Key) interface{} {
				if v, e := strconv.ParseUint(k.Base[0].Value, 10, 64); e != nil {
					panic(e)
				} else {
					return v
				}
			},
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
		},
	)
)
