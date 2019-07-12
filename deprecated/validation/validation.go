package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	model "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersona(bs []byte, persona *model.Persona) *fabric.Response {
	if res := validatePersonaJSON(bs); res != nil {
		return res
	}
	if err := json.Unmarshal(bs, persona); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Unmarshaling Persona: %s", err))
	}
	return &fabric.Response{}
}

// TODO refactorizar en funcionalidades de al menos 2 capas: unmarshalling de personas y respuesta de servicio
func ArgToPersonas(bs []byte, personas *[]model.Persona) *fabric.Response {
	if res := validatePersonaListJSON(bs); res != nil {
		return res
	}
	if err := json.Unmarshal(bs, personas); err != nil {
		return fabric.SystemErrorResponse("Unmarshalling Personas: " + err.Error())
	}
	return &fabric.Response{}
}

var personaSchema = validator.MustLoadSchema("persona")

func validatePersonaJSON(bs []byte) *fabric.Response {
	res, err := validator.Validate(personaSchema, bs)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de persona mal formado: %s", err))
	}
	if !res.Valid() {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de persona inválido: %s", report(res)))
	}
	return nil
}

var personaListSchema = validator.MustLoadSchema("persona-list")

func validatePersonaListJSON(bs []byte) *fabric.Response {
	res, err := validator.Validate(personaListSchema, bs)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de lista de personas mal formado: %s", err))
	}
	if !res.Valid() {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de lista de persona inválido: %s", report(res)))
	}
	return nil
}

func report(res *validator.ValidationResult) string {
	ss := []string{}
	for _, err := range res.Errors {
		ss = append(ss, fmt.Sprintf("%s: %s", err.Field, err.Description))
	}
	return fmt.Sprintf("%d problemas encontrados:\n%s", len(ss), strings.Join(ss, "\n"))
}
