package personas

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var (
	personasSchemaLoader gojsonschema.JSONLoader
	personaSchemaLoader  gojsonschema.JSONLoader
)

func init() {
	validateJSON("PersonasSchema", personasSchema)
	validateJSON("PersonaSchema", personaSchema)
	validateJSON("ImpuestoSchema", impuestoSchema)
	validateJSON("ActividadSchema", actividadSchema)
	validateJSON("DomicilioSchema", domicilioSchema)
	validateJSON("TelefonoSchema", telefonoSchema)
	personasSchemaLoader = gojsonschema.NewStringLoader(personasSchema)
	personaSchemaLoader = gojsonschema.NewStringLoader(personaSchema)
}

// validateJSON valida un schema y lanza panic porque solo debe fallar en tiempo de desarrollo
// SOLO PARA SER INVOCADADA DESDE init()
func validateJSON(name string, schema string) {
	// esta función lanza panic porque tiene la intención de fallar únicamente
	// en tiempo de desarrollo ya que los schemas están embebidos en
	// código fuente
	_, err := gojsonschema.NewStringLoader(schema).LoadJSON()
	if err != nil {
		panic(fmt.Sprintf("JSON mal formado en el schema %q: %v", name, err))
	}
	if !json.Valid([]byte(schema)) {
		panic(fmt.Sprintf("JSON mal formado en el schema %q", name))
	}
}

const personasSchema = `
{
  "description": "A representation of a person",
  "type": "object",
  "properties": {
    "personas": {
      "type": "array",
      "items": ` + personaSchema + `
    }
  }
}
`

const personaSchema = `
{
  "description": "A representation of a person",
  "type": "object",
  "required": [
    "cuit"
  ],
  "properties": {
    "cuit": {
      "type": "integer"
    },
    "nombre": {
      "type": "string"
    },
    "apellido": {
      "type": "string"
    },
    "materno": {
      "type": "string"
    },
    "razonSocial": {
      "type": "string"
    },
    "tipo": {
      "type": "string"
    },
    "estado": {
      "type": "string"
    },
    "formaJuridica": {
      "type": "integer"
    },
    "tipoDoc": {
      "type": "integer"
    },
    "doc": {
      "type": "string"
    },
    "sexo": {
      "type": "string"
    },
    "mesCierre": {
      "type": "integer"
    },
    "nacimiento": {
      "type": "string"
    },
    "fallecimiento": {
      "type": "string"
    },
    "inscripcion": {
      "type": "string"
    },
    "pais": {
      "type": "string"
    },
    "nuevaCuit": {
      "type": "integer"
    },
    "ch": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": [
          "nombre",
          "apellido",
          "materno",
          "razonSocial",
          "tipo",
          "estado",
          "formaJuridica",
          "tipoDoc",
          "doc",
          "sexo",
          "mesCierre",
          "nacimiento",
          "fallecimiento",
          "inscripcion",
          "pais",
          "nuevaCuit"
        ]
      }
    },
    "ds": {
      "type": "string"
    },
    "impuestos": {
      "type": "array",
      "items": ` + impuestoSchema + `
    },
    "actividades": {
      "type": "array",
      "items": ` + actividadSchema + `
    },
    "domicilios": {
      "type": "array",
      "items": ` + domicilioSchema + `
    },
    "telefonos": {
      "type": "array",
      "items": ` + telefonoSchema + `
    }
  }
}
`

const impuestoSchema = `
{
  "type": "object",
  "required": [
    "impuesto",
    "periodo",
    "estado"
  ],
  "properties": {
    "impuesto": {
      "type": "integer"
    },
    "inscripcion": {
      "type": "string"
    },
    "periodo": {
      "type": "integer"
    },
    "estado": {
      "type": "string"
    },
    "motivo": {
      "type": "string"
    },
    "dia": {
      "type": "integer"
    },
    "ds": {
      "type": "string",
      "minLength": 10,
      "maxLength": 10
    }
  }
}
`

const actividadSchema = `
{
  "type": "object",
  "properties": {
    "id": {
      "type": "integer"
    },
    "nomenclador": {
      "type": "integer"
    },
    "orden": {
      "type": "integer"
    },
    "estado": {
      "type": "string",
      "minLength": 2,
      "maxLength": 2
    },
    "periodo": {
      "type": "integer"
    },
    "ds": {
      "type": "string",
      "minLength": 10,
      "maxLength": 10
    }
  }
}
`

const domicilioSchema = `
{
  "type": "object",
  "required": [
    "orden"
  ],
  "properties": {
    "tipo": {
      "type": "integer"
    },
    "orden": {
      "type": "integer"
    },
    "estado": {
      "type": "integer"
    },
    "nomenclador": {
      "type": "integer"
    },
    "codPostal": {
      "type": "string"
    },
    "provincia": {
      "type": "integer",
      "minimum": 0,
      "maximum": 24
    },
    "localidad": {
      "type": "string"
    },
    "calle": {
      "type": "string"
    },
    "numero": {
      "type": "string"
    },
    "piso": {
      "type": "string"
    },
    "oficina": {
      "type": "string"
    },
    "ds": {
      "type": "string",
      "minLength": 10,
      "maxLength": 10
    }
  }
}
`

const telefonoSchema = `
{
  "type": "object",
  "required": [
    "numero"
  ],
  "properties": {
    "numero": {
      "type": "string"
    },
    "idEstadoTelefono": {
      "type": "string"
    },
    "ds": {
      "type": "string",
      "minLength": 10,
      "maxLength": 10
    }
  }
}
`
