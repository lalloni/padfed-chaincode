package personas_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/test"
)

func TestPutPersona(t *testing.T) {
	stub := test.SetInitTests(t)

	// Valid
	res := test.PutPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", res.Message)
		t.FailNow()
	}

	// Invalid cuit
	res = test.PutPersona(t, stub, 1)
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con un cuit invalido debe dar error")
		t.FailNow()
	}

	// distinct cuits
	var personaJSON = test.GetPersonaJSON(20255438795)
	res = stub.MockInvoke("1", [][]byte{[]byte("putPersona"), []byte("30679638943"), []byte(personaJSON)})
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con cuits distintos debe dar error")
		t.FailNow()
	}
}

const pJSON = `
{
  "personas": [
    {
      "cuit": 20066675573,
      "apellido": "GES",
      "nombre": "THOMAS MICHAEL",
      "tipo": "F",
      "estado": "A",
      "tipoDoc": 1,
      "doc": "6667557",
      "sexo": "M",
      "nacimiento": "1928-11-17",
      "impuestos": [
        {
          "impuesto": 20,
          "estado": "BD",
          "periodo": 199901
        },
        {
          "impuesto": 5243,
          "estado": "BD",
          "periodo": 200907
        },
        {
          "impuesto": 21,
          "estado": "BD",
          "periodo": 200907
        },
        {
          "impuesto": 5244,
          "estado": "AC",
          "periodo": 199807
        }
      ],
      "categorias": [
        {
          "idCategoria": "11",
          "estado": "BD",
          "impuesto": 20,
          "periodo": 200907
        },
        {
          "idCategoria": "11",
          "estado": "BD",
          "impuesto": 23,
          "periodo": 200907
        }
      ],
      "actividades": [
        {
          "nomenclador": 883,
          "id": 692000,
          "orden": 1,
          "estado": "AC",
          "periodo": 201311
        }
      ],
      "domicilios": [
        {
          "tipo": 1,
          "orden": 1,
          "estado": 2,
          "nomenclador": 3541,
          "codPostal": "5891",
          "provincia": 3,
          "localidad": "VILLA CURA BROCHERO",
          "calle": "HIPOLITO IRIGOYEN",
          "numero": "57"
        },
        {
          "tipo": 2,
          "orden": 1,
          "estado": 9,
          "nomenclador": 3541,
          "codPostal": "5891",
          "provincia": 3,
          "localidad": "VILLA CURA BROCHERO",
          "calle": "SAN MARTIN ESQ IRIGO",
          "numero": "8"
        }
      ]
    },
    {
      "cuit": 20066758193,
      "apellido": "RACCONTARE",
      "nombre": "GUSTAVO FABIAN",
      "tipo": "F",
      "estado": "A",
      "tipoDoc": 1,
      "doc": "6675819",
      "sexo": "M",
      "nacimiento": "1933-01-22",
      "impuestos": [
        {
          "impuesto": 20,
          "estado": "AC",
          "periodo": 190101
        },
        {
          "impuesto": 30,
          "estado": "AC",
          "periodo": 200408
        },
        {
          "impuesto": 32,
          "estado": "BD",
          "periodo": 200408
        },
        {
          "impuesto": 5244,
          "estado": "AC",
          "periodo": 199105
        },
        {
          "impuesto": 301,
          "estado": "AC",
          "periodo": 199407
        },
        {
          "impuesto": 5100,
          "estado": "AC",
          "periodo": 196501
        }
      ],
      "categorias": [
        {
          "idCategoria": "501",
          "estado": "AC",
          "impuesto": 5100,
          "periodo": 200703
        }
      ],
      "actividadList": [],
      "domicilios": [
        {
          "tipo": 1,
          "orden": 1,
          "estado": 6,
          "nomenclador": 6024,
          "codPostal": "3315",
          "provincia": 19,
          "localidad": "LEANDRO N. ALEM",
          "calle": "RIVADAVIA",
          "numero": "572"
        },
        {
          "tipo": 2,
          "orden": 1,
          "estado": 1,
          "nomenclador": 6024,
          "codPostal": "3315",
          "provincia": 19,
          "localidad": "LEANDRO N. ALEM",
          "calle": "URUGUAY",
          "numero": "287"
        }
      ]
    }
  ]
}
`

func TestPutPersonas(t *testing.T) {
	stub := test.SetInitTests(t)
	b := bytes.Buffer{}
	if err := json.Compact(&b, []byte(pJSON)); err != nil {
		t.Fatalf("compacting json: %v", err)
	}
	res := stub.MockInvoke("1", [][]byte{[]byte("putPersonas"), b.Bytes()})
	if res.Status != shim.OK {
		fmt.Println("putPersonas", res.Message)
		t.FailNow()
	} else {
		fmt.Println("putPersonas Ok!!!!")
	}

}
