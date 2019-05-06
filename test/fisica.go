package test

import (
	"bytes"
	"text/template"
)

func GetPersonaFísica(cuit uint64) []byte {
	buf := &bytes.Buffer{}
	err := fist.Execute(buf, map[string]interface{}{"CUIT": cuit})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()

}

var fist = template.Must(template.New("personafisica").Parse(`
{
    "id": {{ .CUIT }},
    "persona":{
        "id": {{ .CUIT }},
        "tipoid":"C",
        "tipo":"F",
        "estado":"A",
        "nombre":"XXXXXXX XXXX",
        "apellido":"XXXXXX",
        "materno":"XXXXXX",
        "sexo":"M",
        "nacimiento":"1961-12-30",
        "documento":{
            "tipo":96,
            "numero":"XXXXXXXX"
        }
    },
    "impuestos":{
        "30":{
            "impuesto":30,
            "periodo":198903,
            "estado":"AC",
            "dia":1,
            "motivo":44,
            "inscripcion":"1989-03-01",
            "ds":"2003-06-07"
        },
        "301":{
            "impuesto":301,
            "periodo":198010,
            "estado":"AC",
            "dia":1,
            "motivo":44,
            "inscripcion":"1980-10-01",
            "ds":"2003-06-07"
        }
    },
    "actividades":{
        "883-120091":{
            "actividad":"883-120091",
            "orden":1,
            "desde":201311,
            "ds":"2014-10-02"
        },
        "883-120099":{
            "actividad":"883-120099",
            "orden":2,
            "desde":201311,
            "ds":"2014-10-02"
        }
    },
    "domicilios":{
        "1.1":{
            "tipo":1,
            "orden":1,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXXXXXXX",
            "numero":26950,
            "ds":"2004-02-04"
        },
        "2.1":{
            "tipo":2,
            "orden":1,
            "estado":2,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "calle":"XX XXXXXXXXX",
            "numero":26950,
            "ds":"2003-06-07"
        },
        "3.1":{
            "tipo":3,
            "orden":1,
            "estado":6,
            "provincia":7,
            "localidad":"MAIPU",
            "cp":"5515",
            "nomenclador":"5667",
            "calle":"XXXXXXXXX XX¿X",
            "numero":2900,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2004-10-06"
        },
        "3.2":{
            "tipo":3,
            "orden":2,
            "estado":6,
            "provincia":16,
            "localidad":"RESISTENCIA",
            "cp":"3500",
            "nomenclador":"3983",
            "calle":"XXXX XXX XX XX XXXX",
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2004-10-06"
        },
        "3.3":{
            "tipo":3,
            "orden":3,
            "estado":6,
            "provincia":14,
            "localidad":"BANDA DEL RIO SALI",
            "cp":"4109",
            "nomenclador":"10482",
            "calle":"XXX XXXXXXXX X XXXX XXX XX X.X",
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2004-10-06"
        },
        "3.4":{
            "tipo":3,
            "orden":4,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXX XXX",
            "numero":1551,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.5":{
            "tipo":3,
            "orden":5,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXXXXXXX",
            "numero":26950,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXXXXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.6":{
            "tipo":3,
            "orden":6,
            "estado":6,
            "provincia":4,
            "localidad":"GOYA",
            "cp":"3450",
            "nomenclador":"3850",
            "calle":"XXXX XXXXXXX",
            "numero":93,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.7":{
            "tipo":3,
            "orden":7,
            "estado":6,
            "provincia":19,
            "localidad":"LEANDRO N. ALEM",
            "cp":"3315",
            "nomenclador":"6024",
            "calle":"XXXX X XXXXXXXX",
            "numero":603,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.8":{
            "tipo":3,
            "orden":8,
            "estado":6,
            "provincia":9,
            "localidad":"ROSARIO DE LERMA",
            "cp":"4405",
            "nomenclador":"6505",
            "calle":"XXX XXXXX",
            "numero":370,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.9":{
            "tipo":3,
            "orden":9,
            "estado":6,
            "provincia":6,
            "localidad":"PERICO",
            "cp":"4608",
            "nomenclador":"5123",
            "calle":"XX XXXXXXXX",
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.10":{
            "tipo":3,
            "orden":10,
            "estado":6,
            "provincia":14,
            "localidad":"LA COCHA",
            "cp":"4162",
            "nomenclador":"10357",
            "calle":"XXXX XXX XX XX XXX",
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-01-07"
        },
        "3.11":{
            "tipo":3,
            "orden":11,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXX XXX",
            "numero":1751,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.12":{
            "tipo":3,
            "orden":12,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXXXXXXX",
            "numero":26950,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.13":{
            "tipo":3,
            "orden":13,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XXXXX",
            "numero":55,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXX"
            },
            "ds":"2008-02-19"
        },
        "3.14":{
            "tipo":3,
            "orden":14,
            "estado":2,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "calle":"XXXX X XX XX X/X",
            "adicional":{
                "tipo":5,
                "dato":"XXX,XXXXX/XXX. XXXX"
            },
            "ds":"2008-02-19"
        },
        "3.15":{
            "tipo":3,
            "orden":15,
            "estado":2,
            "provincia":4,
            "localidad":"GOYA",
            "cp":"3450",
            "calle":"XXXXXXXX XXXXXXXXX XXX",
            "adicional":{
                "tipo":5,
                "dato":"XXX,XXXXX/XXX. XXXX"
            },
            "ds":"2008-02-19"
        },
        "3.16":{
            "tipo":3,
            "orden":16,
            "estado":6,
            "provincia":1,
            "localidad":"MAR DEL PLATA SUR",
            "cp":"7600",
            "nomenclador":"1345",
            "calle":"XXXXX XXXX",
            "numero":2964,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.17":{
            "tipo":3,
            "orden":17,
            "estado":6,
            "provincia":3,
            "localidad":"CIUDAD DE CORDOBA SUR (NO DISTRIBUIDOS)",
            "cp":"5008",
            "nomenclador":"2305",
            "calle":"XXXXXX",
            "numero":2792,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.18":{
            "tipo":3,
            "orden":18,
            "estado":6,
            "provincia":22,
            "localidad":"CIPOLLETTI",
            "cp":"8324",
            "nomenclador":"6253",
            "calle":"XXXXXXXX",
            "numero":50,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.19":{
            "tipo":3,
            "orden":19,
            "estado":6,
            "provincia":0,
            "cp":"1003",
            "nomenclador":"401003020",
            "calle":"XXXX XXXXXXX X. XX.",
            "numero":466,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2008-02-19"
        },
        "3.20":{
            "tipo":3,
            "orden":20,
            "estado":6,
            "provincia":1,
            "localidad":"MERLO",
            "cp":"1722",
            "nomenclador":"72",
            "calle":"XX XXXXXXXXX",
            "numero":26950,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXX"
            },
            "ds":"2008-04-24"
        },
        "3.21":{
            "tipo":3,
            "orden":21,
            "estado":6,
            "provincia":1,
            "localidad":"JOSE CLEMENTE PAZ",
            "cp":"1665",
            "nomenclador":"368",
            "calle":"XXXXXXX XXXXXXX X/XXX",
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXX"
            },
            "ds":"2008-04-24"
        },
        "3.22":{
            "tipo":3,
            "orden":22,
            "estado":6,
            "provincia":12,
            "localidad":"ROSARIO SUD",
            "cp":"2000",
            "nomenclador":"7825",
            "calle":"XXXXXXX",
            "numero":4450,
            "adicional":{
                "tipo":5,
                "dato":"XXXXXXXX XXXXXXXXX"
            },
            "ds":"2008-02-19"
        }
    },
    "telefonos":{
        "1":{
            "orden":1,
            "pais":200,
            "area":220,
            "numero":9999999,
            "tipo":6,
            "linea":1,
            "ds":"2013-12-16"
        }
    }
}
`))
