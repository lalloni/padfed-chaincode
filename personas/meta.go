package personas

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var Persona = meta.MustPrepare(meta.Composite{
	Name:       "persona",
	Creator:    personaCreator,
	Identifier: func(v interface{}) interface{} { return v.(*model.Persona).ID },
	Keyer:      func(id interface{}) *key.Key { return key.Based("per", strconv.FormatUint(id.(uint64), 10)) },
	Singletons: []meta.Singleton{
		{
			Tag:     "per",
			Creator: func() interface{} { return &model.PersonaBasica{} },
			Getter:  func(v interface{}) interface{} { return v.(*model.Persona).Persona },
			Setter:  func(v interface{}, w interface{}) { v.(*model.Persona).Persona = w.(*model.PersonaBasica) },
		},
	},
	Collections: []meta.Collection{
		{
			Tag:        "act",
			Creator:    func() interface{} { return &model.PersonaActividad{} },
			Collector:  personaActividadCollect,
			Enumerator: personaActividadEnum,
		},
		{
			Tag:        "imp",
			Creator:    func() interface{} { return &model.PersonaImpuesto{} },
			Collector:  personaImpuestoCollect,
			Enumerator: personaImpuestoEnum,
		},
		{
			Tag:        "dom",
			Creator:    func() interface{} { return &model.PersonaDomicilio{} },
			Collector:  personaDomicilioCollect,
			Enumerator: personaDomicilioEnum,
		},
	},
})

func personaCreator() interface{} {
	return &model.Persona{
		Actividades:    map[string]*model.PersonaActividad{},
		Impuestos:      map[string]*model.PersonaImpuesto{},
		Domicilios:     map[string]*model.PersonaDomicilio{},
		Telefonos:      map[string]*model.PersonaTelefono{},
		Jurisdicciones: map[string]*model.PersonaJurisdiccion{},
		Emails:         map[string]*model.PersonaEmail{},
		Archivos:       map[string]*model.PersonaArchivo{},
		Categorias:     map[string]*model.PersonaCategoria{},
		Etiquetas:      map[string]*model.PersonaEtiqueta{},
		Contribuciones: map[string]*model.PersonaContribucion{},
		Relaciones:     map[string]*model.PersonaRelacion{},
	}
}

func personaImpuestoCollect(v interface{}, i meta.Item) {
	v.(*model.Persona).Impuestos[i.Identifier] = i.Value.(*model.PersonaImpuesto)
}

func personaImpuestoEnum(v interface{}, items *[]meta.Item) {
	for id, w := range v.(*model.Persona).Impuestos {
		*items = append(*items, meta.NewItem(id, w))
	}
}

func personaActividadCollect(v interface{}, i meta.Item) {
	v.(*model.Persona).Actividades[i.Identifier] = i.Value.(*model.PersonaActividad)
}

func personaActividadEnum(v interface{}, items *[]meta.Item) {
	for id, w := range v.(*model.Persona).Actividades {
		*items = append(*items, meta.NewItem(id, w))
	}
}

func personaDomicilioCollect(v interface{}, i meta.Item) {
	v.(*model.Persona).Domicilios[i.Identifier] = i.Value.(*model.PersonaDomicilio)
}

func personaDomicilioEnum(v interface{}, items *[]meta.Item) {
	for id, w := range v.(*model.Persona).Domicilios {
		*items = append(*items, meta.NewItem(id, w))
	}
}
