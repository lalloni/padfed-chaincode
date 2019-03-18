package personas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/lalloni/afip/cuit"
	"github.com/xeipuuv/gojsonschema"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func applySchema(bs []byte, schemaLoader gojsonschema.JSONLoader) *fabric.Response {
	documentLoader := gojsonschema.NewBytesLoader(bs)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON mal formado: %s", err))
	}
	if !result.Valid() {
		report := &strings.Builder{}
		for _, err := range result.Errors() {
			report.WriteString(err.String())
			report.WriteString("\n")
		}
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON no cumple con el esquema: %s", report.String()))
	}
	return nil
}

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersona(bs []byte, persona *model.Persona) *fabric.Response {
	if res := applySchema(bs, personaSchemaLoader); res != nil {
		return res
	}
	if err := json.Unmarshal(bs, persona); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Unmarshaling Persona: %s", err))
	}
	return ValidatePersona(persona)
}

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersonas(bs []byte, personas *model.Personas) *fabric.Response {
	if res := applySchema(bs, personasSchemaLoader); res != nil {
		return res
	}
	if err := json.Unmarshal(bs, &personas); err != nil {
		return fabric.SystemErrorResponse("Unmarshalling Personas: " + err.Error())
	}
	for _, p := range personas.Personas {
		if err := ValidatePersona(p); !err.IsOK() {
			return err
		}
	}
	return &fabric.Response{}
}

func fmtInvalidField(f string, v, err interface{}) string {
	s := fmt.Sprintf("atributo '%s' con valor [%v] inválido", f, v)
	if err != nil {
		s = fmt.Sprintf("%s: %v", s, err)
	}
	return s
}

// TODO refactorizar en funcionalidades de al menos 2 capas: negocio de validación y respuesta de servicio
func ValidatePersona(persona *model.Persona) *fabric.Response {
	if !cuit.IsValid(persona.CUIT) {
		return fabric.ClientErrorResponse(fmtInvalidField("cuit", persona.CUIT, nil))
	}
	if err := helpers.ValidateDate(persona.Nacimiento); err != nil {
		return fabric.ClientErrorResponse(fmtInvalidField("nacimiento", persona.Nacimiento, err))
	}
	if err := helpers.ValidateDate(persona.Inscripcion); err != nil {
		return fabric.ClientErrorResponse(fmtInvalidField("inscripcion", persona.Inscripcion, err))
	}
	if err := helpers.ValidateDate(persona.FechaCierre); err != nil {
		return fabric.ClientErrorResponse(fmtInvalidField("fechaCierre", persona.FechaCierre, err))
	}
	if err := helpers.ValidateDate(persona.Fallecimiento); err != nil {
		return fabric.ClientErrorResponse(fmtInvalidField("fallecimiento", persona.Fallecimiento, err))
	}
	if err := helpers.ValidateDate(persona.DS); err != nil {
		return fabric.ClientErrorResponse(fmtInvalidField("ds", persona.DS, err))
	}
	return &fabric.Response{}
}

func GetPersonaKey(persona *model.Persona) string {
	cuitStr := strconv.FormatUint(persona.CUIT, 10)
	return "PER_" + cuitStr
}
