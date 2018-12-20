package main

const ImpuestoSchema = `{
    "type": "object",
    "required": ["idImpuesto"],
    "properties": {
        "idImpuesto": {
            "type": "number"
        },
        "idOrg": {
            "type": "number"
        },
        "fechaInscripcion": {
            "type": "string"
        },
        "periodo": {
            "type": "number"
        },
        "estado": {
            "type": "string"
        },
        "idTxc": {
            "type": "number"
        }
    }
}`

const ActividadSchema = `{
    "type": "object",
    "required": ["idActividad"],
    "properties": {
        "codNomenclador": {
            "type": "number"
        },
        "idActividad": {
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
    "required": ["idTipoDomicilio"],
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
    "required": ["cuit", "estadoCuit", "tipoPersona"],
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
        "tipoPersona": {
            "type": "string"
        },
        "estadoCuit": {
            "type": "string"
        },
        "idFormaJuridica": {
            "type": "number"
        },
        "tipoDoc": {
            "type": "number"
        },
        "documento": {
            "type": "string"
        },
        "sexo": {
            "type": "string"
        },
        "mesCierre": {
            "type": "number"
        },
        "fechaNacimiento": {
            "type": "string"
        },
        "fechaFallecimiento": {
            "type": "string"
        },
        "fechaInscripcion": {
            "type": "string"
        },
        "fechaCierre": {
            "type": "string"
        },
        "nuevaCuit": {
            "type": "number"
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
