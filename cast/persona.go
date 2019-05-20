package cast

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var Persona = meta.MustPrepare(meta.Composite{
	Name:            "persona",
	KeyBaseName:     "per",
	IdentifierField: "ID",
	IdentifierKey:   Uint64Key("per"),
	KeyIdentifier:   Uint64Identifier(0),
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

func Uint64Key(name string) meta.KeyFunc {
	return func(id interface{}) (*key.Key, error) {
		return key.NewBase(name, strconv.FormatUint(id.(uint64), 10)), nil
	}
}

func Uint64Identifier(seg int) meta.ValFunc {
	return func(k *key.Key) (interface{}, error) {
		return strconv.ParseUint(k.Base[seg].Value, 10, 64)
	}
}
