package personas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
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
func ArgToPersonas(bs []byte, personas *model.Personas) *fabric.Response {
	if res := validatePersonaListJSON(bs); res != nil {
		return res
	}
	if err := json.Unmarshal(bs, personas); err != nil {
		return fabric.SystemErrorResponse("Unmarshalling Personas: " + err.Error())
	}
	return &fabric.Response{}
}

func validatePersonaJSON(bs []byte) *fabric.Response {
	res, err := validator.ValidatePersonaJSON(bs)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de persona mal formado: %s", err))
	}
	if !res.Valid() {
		return fabric.ClientErrorResponse(fmt.Sprintf("JSON de persona inválido: %s", report(res)))
	}
	return nil
}

func validatePersonaListJSON(bs []byte) *fabric.Response {
	res, err := validator.ValidatePersonaListJSON(bs)
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

func GetPersonaKey(persona *model.Persona) string {
	return "PER_" + strconv.FormatUint(persona.ID, 10)
}
