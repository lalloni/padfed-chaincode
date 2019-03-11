package personas

import (
	"encoding/json"
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"

	"github.com/lalloni/afip/cuit"
	"github.com/xeipuuv/gojsonschema"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersona(personaAsBytes []byte, persona *model.Persona) *fabric.Response {
	documentLoader := gojsonschema.NewStringLoader(string(personaAsBytes))
	result, err := gojsonschema.Validate(personaSchemaLoader, documentLoader)
	if err != nil {
		return fabric.ClientErrorResponse("JSON schema invalido: " + err.Error() + " - " + string(personaAsBytes))
	}

	if !result.Valid() {
		var errosStr string
		for _, desc := range result.Errors() {
			errosStr += desc.Description() + ". "
		}
		return fabric.ClientErrorResponse("JSON no cumple con el esquema: " + errosStr)
	}
	err = json.Unmarshal(personaAsBytes, persona)
	if err != nil {
		return fabric.SystemErrorResponse("JSON invalido: " + err.Error())
	}
	return ValidatePersona(persona)
}

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersonas(personasAsBytes []byte, personas *model.Personas) *fabric.Response {
	documentLoader := gojsonschema.NewStringLoader(string(personasAsBytes))
	result, err := gojsonschema.Validate(personasSchemaLoader, documentLoader)
	if err != nil {
		return fabric.ClientErrorResponse("JSON schema invalido: " + err.Error() + " - " + string(personasAsBytes))
	}
	if !result.Valid() {
		var errosStr string
		for _, desc := range result.Errors() {
			errosStr += desc.Description() + ". "
		}
		return fabric.ClientErrorResponse("JSON no cumple con el esquema: " + errosStr)
	}
	err = json.Unmarshal(personasAsBytes, &personas)
	if err != nil {
		return fabric.SystemErrorResponse("JSON invalido: " + err.Error())
	}

	for _, p := range personas.Personas {
		err := ValidatePersona(p)
		if !err.IsOK() {
			return err
		}
	}
	return &fabric.Response{}
}

// TODO refactorizar en funcionalidades de al menos 2 capas: negocio de validación y respuesta de servicio
func ValidatePersona(persona *model.Persona) *fabric.Response {
	var err fabric.Response
	cuitStr := strconv.FormatUint(persona.CUIT, 10)
	if !cuit.IsValid(persona.CUIT) {
		return fabric.ClientErrorResponse("cuit [" + cuitStr + "] invalida")
	}
	if err := helpers.ValidateDate(persona.Nacimiento); err != nil {
		return fabric.ClientErrorResponse("nacimiento [" + persona.Nacimiento + "] invalido: " + err.Error())
	}
	if err := helpers.ValidateDate(persona.Inscripcion); err != nil {
		return fabric.ClientErrorResponse("inscripcion [" + persona.Inscripcion + "] invalida: " + err.Error())
	}
	if err := helpers.ValidateDate(persona.FechaCierre); err != nil {
		return fabric.ClientErrorResponse("fechaCierre [" + persona.FechaCierre + "] invalida: " + err.Error())
	}
	if err := helpers.ValidateDate(persona.Fallecimiento); err != nil {
		return fabric.ClientErrorResponse("fallecimiento [" + persona.Fallecimiento + "] invalido: " + err.Error())
	}
	if err := helpers.ValidateDate(persona.DS); err != nil {
		return fabric.ClientErrorResponse("ds [" + persona.DS + "] invalido: " + err.Error())
	}
	return &err
}

func GetPersonaKey(persona *model.Persona) string {
	cuitStr := strconv.FormatUint(persona.CUIT, 10)
	return "PER_" + cuitStr
}
