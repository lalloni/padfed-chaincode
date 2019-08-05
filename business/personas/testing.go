package personas

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/lalloni/afip/cuit"
	words "github.com/moul/number-to-words"

	"github.com/lalloni/fabrikit/chaincode/test"
)

func SummaryPersonasID(pers []Persona) (index map[uint64]Persona, ids []uint64) {
	index = map[uint64]Persona{}
	ids = []uint64{}
	for _, per := range pers {
		index[per.ID] = per
	}
	for id := range index {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return
}

func RandomPersonas(q int, rnd *rand.Rand) []Persona {
	if rnd == nil {
		rnd = test.NewTimeRand()
	}
	check := map[uint64]bool{}
	pers := []Persona{}
	for n := 0; n < q; n++ {
		id := cuit.Random(rnd)
		for check[id] { // don't repeat
			id = cuit.Random(rnd)
		}
		check[id] = true
		per := Persona{
			ID: id,
			Persona: &Basica{
				ID:     id,
				TipoID: "C",
				Estado: "A",
			},
			Emails: map[string]*Email{
				"1": {Direccion: "foo@bar.com", Orden: 1, Tipo: 1, Estado: 1},
			},
			Archivos: map[string]*Archivo{
				"1": {Descripcion: "Archivo X", Orden: 1, Tipo: 10},
			},
			Etiquetas: map[string]*Etiqueta{
				"1":  {Etiqueta: 1, Periodo: 20101001, Estado: "AC"},
				"10": {Etiqueta: 10, Periodo: 20101001, Estado: "AC"},
			},
		}
		switch cuit.TipoPersonaCUIT(id) {
		case cuit.PersonaFísica:
			per.Persona.Tipo = "F"
			per.Persona.Nombre = strings.ToUpper(words.IntegerToEsEs(int(id % 10)))
			per.Persona.Apellido = strings.ToUpper(words.IntegerToEsEs(int(id/10) % 100))
		case cuit.PersonaJurídica:
			per.Persona.Tipo = "J"
			per.Persona.RazonSocial = strings.ToUpper(words.IntegerToEsEs(int(id % 100)))
		}
		pers = append(pers, per)
	}
	return pers
}
