package main

const ImpuestoSchema = `{
    "type": "object",
    "required": ["impuesto", "periodo", "estado"],
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
            "minLength":10,
            "maxLength":10
        }
    }
}`

const ActividadSchema = `{
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
            "type": "string"
            "minLength":2,
            "maxLength":2
        },
        "periodo": {
            "type": "integer"
        },
        "ds": {
            "type": "string",
            "minLength":10,
            "maxLength":10
        }
    }
}`

const DomicilioSchema = `{
    "type": "object",
    "required": ["orden"],
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
            "minimum":0,
            "maximum":24
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
            "minLength":10,
            "maxLength":10
        }
    }
}`

const TelefonoSchema = `{
    "type": "object",
    "required": ["numero"],
    "properties": {
        "numero": {
            "type": "string"
        },
        "idEstadoTelefono": {
            "type": "string"
        },
        "ds": {
            "type": "string",
            "minLength":10,
            "maxLength":10
        }
    }
}`

const PersonaSchema = `{
    "description": "A representation of a person",
    "type": "object",
    "required": ["cuit"],
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
                    "nombre", "apellido", "materno", "razonSocial", "tipo", "estado",
                    "formaJuridica", "tipoDoc", "doc", "sexo", "mesCierre", "nacimiento",
                    "fallecimiento", "inscripcion", "pais", "nuevaCuit"
                ]
            }
        },
        "ds": {
            "type": "string"
        },
        "impuestos": {
            "type": "array",
            "items": ` + ImpuestoSchema + `},
        "actividades": {
            "type": "array",
            "items": ` + ActividadSchema + `},
        "domicilios": {
            "type": "array",
            "items": ` + DomicilioSchema + `},
        "telefonos": {
            "type": "array",
            "items": ` + TelefonoSchema + `}
    }
}`

const PersonasSchema = `{
    "description": "A representation of a person",
    "type": "object",
    "properties": {
        "personas": {
            "type": "array",
            "items": ` + PersonaSchema + `}
        }
    }`
