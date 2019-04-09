package personas

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func GetPersonaKey(persona *model.Persona) string {
	return GetPersonaKeyCUIT(persona.ID)
}

func GetPersonaKeyCUIT(cuit uint64) string {
	return "PER_" + strconv.FormatUint(cuit, 10)
}
