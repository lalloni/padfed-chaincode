package test

import (
	"math/rand"
	"strings"
	"time"

	"github.com/lalloni/afip/cuit"
	ntw "github.com/moul/number-to-words"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func NewTimeRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomPersonas(q int, rnd *rand.Rand) []model.Persona {
	if rnd == nil {
		rnd = NewTimeRand()
	}
	pers := []model.Persona{}
	for n := 0; n < q; n++ {
		id := cuit.Random(rnd)
		per := model.Persona{
			ID: id,
			Persona: &model.PersonaBasica{
				ID:     id,
				TipoID: "C",
				Estado: "A",
			},
			Etiquetas: map[string]*model.PersonaEtiqueta{
				"1":  {Etiqueta: 1, Periodo: 20101001, Estado: "AC"},
				"10": {Etiqueta: 10, Periodo: 20101001, Estado: "AC"},
			},
		}
		switch cuit.TipoPersonaCUIT(id) {
		case cuit.PersonaFísica:
			per.Persona.Tipo = "F"
			per.Persona.Nombre = strings.ToUpper(ntw.IntegerToEsEs(int(id % 10)))
			per.Persona.Apellido = strings.ToUpper(ntw.IntegerToEsEs(int(id/10) % 100))
		case cuit.PersonaJurídica:
			per.Persona.Tipo = "J"
			per.Persona.RazonSocial = strings.ToUpper(ntw.IntegerToEsEs(int(id % 100)))
		}
		pers = append(pers, per)
	}
	return pers
}
