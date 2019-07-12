package personas

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/convert"
)

func TestMarshaling(t *testing.T) {

	tests := []struct {
		zero interface{}
		v    interface{}
		bs   []byte
	}{
		{&Domicilio{}, &Domicilio{}, []byte("{}")},
		{&Domicilio{}, samplePersonaDomicilio(), samplePersonaDomicilioBytes()},
		{&Email{}, &Email{}, []byte("{}")},
		{&Email{}, samplePersonaEmail(), samplePersonaEmailBytes()},
		// TODO completar para todos los struct Persona*
	}

	for _, test := range tests {

		test := test

		t.Run(fmt.Sprintf("%T/marshal", test.v), func(t *testing.T) {
			a := assert.New(t)
			bs, err := json.Marshal(test.v)
			a.NoError(err, "marshaling %+v", test.v)
			bs, err = convert.Pretty(bs)
			a.NoError(err, "prettyfing %v", bs)
			a.Equal(string(test.bs), string(bs))
		})

		t.Run(fmt.Sprintf("%T/unmarshal", test.v), func(t *testing.T) {
			a := assert.New(t)
			err := json.Unmarshal(test.bs, &test.zero)
			a.NoError(err, "unmarshalling %q", string(test.bs))
			a.Equal(test.zero, test.v)
		})

	}

}

func samplePersonaEmail() *Email {
	return &Email{
		Direccion: "aaa@nowhere.com",
		Orden:     10,
		Tipo:      10,
		Estado:    1,
		DS:        common.FechaEn(2018, 2, 3),
	}
}

func samplePersonaEmailBytes() []byte {
	return []byte(`{
  "direccion": "aaa@nowhere.com",
  "orden": 10,
  "tipo": 10,
  "estado": 1,
  "ds": "2018-02-03"
}`)
}

func samplePersonaDomicilioBytes() []byte {
	return []byte(`{
  "nombre": "pedro",
  "orden": 10,
  "org": 1,
  "tipo": 1,
  "calle": "desconocida",
  "numero": 20,
  "piso": "segundo",
  "sector": "x",
  "manzana": "45",
  "torre": "no tiene",
  "unidad": "alguna",
  "provincia": 0,
  "localidad": "pérez",
  "cp": "C1428FFR",
  "nomenclador": "10",
  "adicional": {
    "tipo": 99,
    "dato": "mmmm"
  },
  "baja": "2019-01-01",
  "ds": "2019-04-08"
}`)
}

func samplePersonaDomicilio() *Domicilio {
	return &Domicilio{
		Nombre:      "pedro",
		Orden:       10,
		Org:         1,
		Tipo:        1,
		Estado:      0,
		Calle:       "desconocida",
		Numero:      20,
		Piso:        "segundo",
		Sector:      "x",
		Manzana:     "45",
		Torre:       "no tiene",
		Unidad:      "alguna",
		Provincia:   common.ProvinciaCódigo(0),
		Localidad:   "pérez",
		CP:          "C1428FFR",
		Nomenclador: "10",
		Adicional: &Adicional{
			Tipo: 99,
			Dato: "mmmm",
		},
		Baja: common.FechaEn(2019, 1, 1),
		DS:   common.FechaEn(2019, 4, 8),
	}
}
