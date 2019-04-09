package personas

import (
	"strconv"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func GetPersonaKey(persona *model.Persona) string {
	return "PER_" + strconv.FormatUint(persona.ID, 10)
}
