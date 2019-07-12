package personas

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/store"

	persona "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/personas"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/fabric"
)

func SavePersona(stub shim.ChaincodeStubInterface, p *persona.Persona) *fabric.Response {

	if p.Persona != nil && p.ID != p.Persona.ID {
		return fabric.ClientErrorResponse(fmt.Sprintf("El root.id [%d] y root.persona.id [%d] deben ser iguales", p.ID, p.Persona.ID))
	}

	st := store.New(stub)

	if exist, err := st.HasComposite(persona.Schema, p.ID); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Error obteniendo existencia de persona: %v", err))
	} else if !exist && p.Persona == nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Se requiere el atributo persona al crear una persona"))
	}

	if err := st.PutComposite(persona.Schema, p); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Error guardando persona: %v", err))
	}

	return fabric.SuccessResponse("OK", 1)

}
