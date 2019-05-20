package test

import (
	"bytes"
	"text/template"
)

func GetPersonaJurídica(cuit uint64) []byte {
	buf := &bytes.Buffer{}
	err := jurt.Execute(buf, map[string]interface{}{"CUIT": cuit})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

var jurt = template.Must(template.New("personajuridica").Parse(`
{
    "id": {{ .CUIT }},
    "persona": {
        "id": {{ .CUIT }},
        "tipoid": "C",
        "tipo": "J",
        "estado": "A",
        "razonsocial": "THE BIRTH OF MARIA CONCETTA",
        "formajuridica": 1,
        "inscripcion": {
            "registro": 10,
            "numero": 10
        },
        "mescierre": 12,
        "ds": "2019-01-01"
    },
    "impuestos": {
        "30": {
            "impuesto": 30,
            "estado": "AC",
            "periodo": 199912
        },
       "5100": {
            "impuesto": 5100,
            "estado": "AC",
            "periodo": 199605
        },
        "301": {
            "impuesto": 301,
            "estado": "AC",
            "periodo": 199407
        },
        "34": {
            "impuesto": 34,
            "estado": "AC",
            "periodo": 201112
        }
    },
    "actividades": {
        "1": {
            "actividad": "883-941100",
            "orden": 1,
            "hasta": 201504,
            "desde": 201311
        }
    },
    "domicilios": {
        "1": {
            "tipo": 1,
            "orden": 1,
            "org": 900,
            "estado": 11,
            "nomenclador": "401084021",
            "cp": "1084",
            "provincia": 0,
            "localidad": "N/A",
            "calle": "DE MAYO AV.",
            "numero": 568
        },
        "2": {
            "tipo": 2,
            "orden": 1,
            "org": 900,
            "estado": 11,
            "nomenclador": "401084021",
            "cp": "1084",
            "provincia": 0,
            "localidad": "N/A",
            "calle": "DE MAYO AV.",
            "numero": 568
        }
    },
    "telefonos": {
        "1": {
            "numero": 46788554,
            "orden": 1
        }
    },
    "cmsedes": {
        "0": {
            "provincia": 2,
            "desde": "2019-06-23",
            "hasta": "2019-04-14",
            "ds": "2019-04-14"
        },
        "1": {
            "provincia": 1,
            "desde": "2019-06-23",
            "ds":"2019-04-14"
        }
    }
}
`))
