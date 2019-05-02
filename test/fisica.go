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
    "persona": {
        "id": {{ .CUIT }},
        "tipoid": "C",
        "tipo": "F",
        "estado": "A",
        "nombre": "Pepe",
        "apellido": "Sanchez",
        "documento": {
            "tipo": 1,
            "numero":  "{{ .CUIT }}"
        },
        "nacimiento": "1928-11-17",
        "ds": "2019-01-01"
    },
    "impuestos": {
        "30": {
            "impuesto": 30,
            "estado": "AC",
            "periodo": 199912,
            "dia": 20
        },
        "5100": {
            "impuesto": 5100,
            "estado": "AC",
            "periodo": 199605,
            "dia": 12
        },
        "301": {
            "impuesto": 301,
            "estado": "AC",
            "periodo": 199407,
            "dia": 12
        },
        "34": {
            "impuesto": 34,
            "estado": "AC",
            "periodo": 201112,
            "dia": 12
        }
    },
    "actividades": {
        "883-941100-1":{
            "actividad": "883-941100",
            "orden": 1,
            "hasta": 201504,
            "desde": 201311
        }
    },
    "domicilios": {
        "1":{
            "tipo": 1,
            "orden": 1,
            "estado": 11,
            "nomenclador": 401084021,
            "cp": "1084",
            "provincia": 0,
            "localidad": "",
            "calle": "DE MAYO AV.",
            "numero": 568
        },
        "2":{
            "tipo": 2,
            "orden": 1,
            "estado": 11,
            "nomenclador": 401084021,
            "cp": "1084",
            "provincia": 0,
            "localidad": "",
            "calle": "DE MAYO AV.",
            "numero": 568
        }
    },
    "telefonos": {
        "46788554": {
            "numero": 46788554,
            "orden": 1
        }
    }
}
`))
