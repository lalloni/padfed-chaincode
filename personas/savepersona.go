package personas

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/cast"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
)

func SavePersona(stub shim.ChaincodeStubInterface, p *model.Persona) *fabric.Response {

	if p.Persona != nil && p.ID != p.Persona.ID {
		return fabric.ClientErrorResponse(fmt.Sprintf("El root.id [%d] y root.persona.id [%d] deben ser iguales", p.ID, p.Persona.ID))
	}

	st := store.New(stub)

	if exist, err := st.HasComposite(cast.Persona, p.ID); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Error obteniendo existencia de persona: %v", err))
	} else if !exist && p.Persona == nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("Se requiere el atributo persona al crear una persona"))
	}

	if err := st.PutComposite(cast.Persona, p); err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Error guardando persona: %v", err))
	}

	return fabric.SuccessResponse("OK", 1)

}
