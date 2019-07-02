package persona

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/lalloni/afip/cuit"
	words "github.com/moul/number-to-words"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/test"
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
	pers := []Persona{}
	for n := 0; n < q; n++ {
		id := cuit.Random(rnd)
		per := Persona{
			ID: id,
			Persona: &Basica{
				ID:     id,
				TipoID: "C",
				Estado: "A",
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
