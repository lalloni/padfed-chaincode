package main

const ImpuestoSchema = `{
    "type": "object",
    "required": ["impuesto", "periodo", "estado"],
    "properties": {
        "impuesto": {
            "type": "number"
        },
        "inscripcion": {
            "type": "string"
        },
        "periodo": {
            "type": "number"
        },
        "estado": {
            "type": "string"
        },
        "ds": {
            "type": "string"
        },
        "motivo": {
            "type": "string"
        },
        "dia": {
            "type": "number"
        }
    }
}`

const ActividadSchema = `{
    "type": "object",
    "properties": {
        "id": {
            "type": "number"
        },
        "codNomenclador": {
            "type": "number"
        },
        "orden": {
            "type": "number"
        },
        "estado": {
            "type": "string"
        },
        "periodo": {
            "type": "number"
        }
    }
}`

const DomicilioSchema = `{
    "type": "object",
    "required": ["orden", "idTipoDomicilio"],
    "properties": {
        "idTipoDomicilio": {
            "type": "number"
        },
        "orden": {
            "type": "number"
        },
        "estado": {
            "type": "string"
        },
        "idNomenclador": {
            "type": "string"
        },
        "codPostal": {
            "type": "string"
        },
        "idProvincia": {
            "type": "string"
        },
        "localidad": {
            "type": "string"
        },
        "calle": {
            "type": "string"
        },
        "numero": {
            "type": "string"
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
        }
    }
}`

const PersonaSchema = `{
    "description": "A representation of a person, company, organization, or place",
    "type": "object",
    "required": ["cuit"],
    "properties": {
        "cuit": {
            "type": "number"
        },
        "nombre": {
            "type": "string"
        },
        "apellido": {
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
            "type": "number"
        },
        "tipoDoc": {
            "type": "number"
        },
        "doc": {
            "type": "string"
        },
        "sexo": {
            "type": "string"
        },
        "mesCierre": {
            "type": "number"
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
        "fechaCierre": {
            "type": "string"
        },
        "nuevaCuit": {
            "type": "number"
        },
        "materno": {
            "type": "string"
        },
        "pais": {
            "type": "string"
        },
        "ch": {
            "type": "string"
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
    "description": "A representation of a person, company, organization, or place",
    "type": "object",
    "properties": {
        "personas": {
            "type": "array",
            "items": ` + PersonaSchema + `}
        }
    }`
