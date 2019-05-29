package personas_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/bitly/go-simplejson"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/test"
)

func TestGetPersona(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)

	cuit := uint64(30679638943)

	personaPUT, err := canonicalize(test.GetPersonaJSON(cuit))
	a.NoError(err)

	res := stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaPUT})
	a.Equal(int32(shim.OK), res.Status, res.Payload)

	res = test.GetPersona(t, stub, cuit)
	a.Equal(int32(shim.OK), res.Status, res.Payload)

	personaGET, err := canonicalize(res.Payload)
	a.NoError(err)

	if !a.Equal(string(personaPUT), string(personaGET)) {
		t.Logf("expected:\n%s\nactual:\n%s", string(personaPUT), string(personaGET))
	}

}

func TestGetPersonaLenientErrors(t *testing.T) {
	a := assert.New(t)

	stub := test.SetInitTests(t)

	cuit := uint64(30679638943)

	personaPUT, err := canonicalize(test.GetPersonaJSON(cuit))
	a.NoError(err)

	res := stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaPUT})
	a.EqualValues(shim.OK, res.Status)

	// modifico y guardo la personabasica inválida
	json, err := simplejson.NewJson(personaPUT)
	a.NoError(err)
	json = json.Get("persona")
	json.Set("pais", "un no número")
	bs, err := json.MarshalJSON()
	a.NoError(err)
	stub.MockTransactionStart("1")
	err = stub.PutState("per:30679638943#per", bs)
	a.NoError(err)
	stub.MockTransactionEnd("1")

	// cargo en modo lenient pidiendo que embeba errores
	res = stub.MockInvoke("1", [][]byte{[]byte("getPersona?lenientread&embederrors"), []byte(strconv.FormatUint(cuit, 10))})
	a.EqualValues(shim.OK, res.Status)

	// controlo que vengan errores
	json, err = simplejson.NewJson(res.Payload)
	a.NoError(err)
	ee, err := json.Get("errors").Array()
	a.NoError(err)
	a.NotEmpty(ee)

}

func canonicalize(bs []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}
	a, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := json.Indent(b, a, "", "  "); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func TestGetPersonaIssue338(t *testing.T) {
	a := assert.New(t)
	stub := test.SetInitTests(t)

	personaPUT, err := canonicalize([]byte(`
        {
            "id": 20293099929,
            "persona": {
                "id": 20293099929,
                "tipoid": "C",
                "tipo": "F",
                "estado": "A",
                "nombre": "XXXXXXX XXXXXX",
                "apellido": "XXXX XXXXX",
                "sexo": "M",
                "nacimiento": "1983-03-18",
                "documento": {
                    "tipo": 96,
                    "numero": "XXXXXXXX"
                }
            },
            "jurisdicciones": {
                "900.3": {
                    "org": 900,
                    "provincia": 3,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                },
                "900.12": {
                    "org": 900,
                    "provincia": 12,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                }
            },
            "cmsedes": {
                "3": {
                    "provincia": 3,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                }
            },
            "impuestos": {
                "11": {
                    "impuesto": 11,
                    "periodo": 200610,
                    "estado": "AC",
                    "dia": 1,
                    "inscripcion": "2005-07-08",
                    "motivo": {
                        "id": 44
                    },
                    "ds": "2006-09-27"
                },
                "30": {
                    "impuesto": 30,
                    "periodo": 200610,
                    "estado": "AC",
                    "dia": 1,
                    "inscripcion": "2005-07-08",
                    "motivo": {
                        "id": 44
                    },
                    "ds": "2006-09-27"
                },
                "308": {
                    "impuesto": 308,
                    "periodo": 200610,
                    "estado": "AC",
                    "dia": 1,
                    "inscripcion": "2005-07-08",
                    "motivo": {
                        "id": 44
                    },
                    "ds": "2006-09-27"
                }
            },
            "categorias": {
                "308.211": {
                    "impuesto": 308,
                    "categoria": 211,
                    "periodo": 200703,
                    "estado": "AC",
                    "ds": "2007-03-15"
                }
            },
            "actividades": {
                "900.900-492221": {
                    "org": 900,
                    "actividad": "900-492221",
                    "orden": 1,
                    "articulo": 9,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                },
                "900.900-492229": {
                    "org": 900,
                    "actividad": "900-492229",
                    "orden": 2,
                    "articulo": 9,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                },
                "900.900-492240": {
                    "org": 900,
                    "actividad": "900-492240",
                    "orden": 3,
                    "articulo": 9,
                    "desde": "2006-10-01",
                    "ds": "2019-05-28"
                }
            },
            "domicilios": {
                "1.1.1": {
                    "org": 1,
                    "tipo": 1,
                    "orden": 1,
                    "estado": 6,
                    "provincia": 3,
                    "localidad": "VILLA DEL ROSARIO",
                    "cp": "5963",
                    "nomenclador": "2423",
                    "calle": "XX XX XXXX",
                    "numero": 886,
                    "ds": "2005-07-22"
                },
                "1.2.1": {
                    "org": 1,
                    "tipo": 2,
                    "orden": 1,
                    "estado": 6,
                    "provincia": 3,
                    "localidad": "VILLA DEL ROSARIO",
                    "cp": "5963",
                    "nomenclador": "2423",
                    "calle": "XX XX XXXX",
                    "numero": 755,
                    "ds": "2005-07-18"
                },
                "900.3.1": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 1,
                    "estado": 6,
                    "provincia": 3,
                    "localidad": "VILLA DEL ROSARIO",
                    "cp": "5963",
                    "calle": "25 DE MAYO",
                    "numero": 480,
                    "ds": "2019-05-28"
                },
                "900.3.2": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 2,
                    "estado": 6,
                    "provincia": 3,
                    "localidad": "VILLA DEL ROSARIO",
                    "cp": "5963",
                    "calle": "25 DE MAYO",
                    "numero": 480,
                    "ds": "2019-05-28"
                },
                "900.3.3": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 3,
                    "estado": 6,
                    "provincia": 3,
                    "localidad": "VILLA DEL ROSARIO",
                    "cp": "5963",
                    "calle": "25 DE MAYO",
                    "numero": 480,
                    "ds": "2019-05-28"
                }
            },
            "domisroles": {
                "900.3.1.1": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 1,
                    "rol": 1,
                    "ds": "2019-05-28"
                },
                "900.3.2.3": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 2,
                    "rol": 3,
                    "ds": "2019-05-28"
                },
                "900.3.3.2": {
                    "org": 900,
                    "tipo": 3,
                    "orden": 3,
                    "rol": 2,
                    "ds": "2019-05-28"
                }
            },
            "telefonos": {
                "1": {
                    "orden": 1,
                    "pais": 200,
                    "area": 3573,
                    "numero": 999999,
                    "tipo": 6,
                    "linea": 1,
                    "ds": "2013-12-16"
                }
            }
        }
	`))
	a.NoError(err)

	res := stub.MockInvoke("1", [][]byte{[]byte("putPersona"), personaPUT})
	a.EqualValues(shim.OK, res.Status)

	res = test.GetPersona(t, stub, 20293099929)
	a.EqualValues(shim.OK, res.Status)

	personaGET, err := canonicalize(res.Payload)
	a.NoError(err)

	if !a.Equal(string(personaPUT), string(personaGET)) {
		t.Logf("expected:\n%s\nactual:\n%s", string(personaPUT), string(personaGET))
	}

}
