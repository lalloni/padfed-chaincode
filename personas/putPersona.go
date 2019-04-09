package personas

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/inscripciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func PutPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se espera 1 (PERSONA).")
	}
	newPersona := &model.Persona{}
	if res := model.ArgToPersona([]byte(args[0]), newPersona); !res.IsOK() {
		return res
	}
	return SavePersona(stub, newPersona)
}

func PutPersonas(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se espera 1 (PERSONAS)")
	}
	newPersonas := &model.PersonaList{}
	if err := model.ArgToPersonas([]byte(args[0]), newPersonas); !err.IsOK() {
		return err
	}

	rows := 0
	for _, p := range newPersonas.Personas {
		p := p
		log.Printf("Grabando persona %d", p.ID)
		res := SavePersona(stub, &p)
		if res.Status != shim.OK {
			res.WrongItem = rows
			return res
		}
		rows++
	}
	log.Println(strconv.Itoa(rows) + " personas processed !!!")
	return fabric.SuccessResponse("Ok", rows)
}

func SavePersona(stub shim.ChaincodeStubInterface, p *model.Persona) *fabric.Response {

	if p.ID != p.Persona.ID {
		return fabric.ClientErrorResponse(fmt.Sprintf("El root.id [%d] y root.persona.id [%d] deben ser iguales", p.ID, p.Persona.ID))
	}

	cuits := strconv.Itoa(int(p.ID))

	key := GetPersonaKey(p)

	exist, err := fabric.KeyExists(stub, cuits)
	if !err.IsOK() {
		return err
	}

	if !exist {
		log.Println("Putting [" + cuits + "]...")
		if err := stub.PutState(cuits, []byte("{}")); err != nil {
			return fabric.SystemErrorResponse("Error putting cuitStr [" + cuits + "]: " + err.Error())
		}
	}

	rows := 0

	if p.Persona != nil {
		bs, err := json.Marshal(p.Persona)
		if err != nil {
			return fabric.SystemErrorResponse("Error marshaling persona [" + key + "]: " + err.Error())
		}
		log.Println("Putting [" + key + "]...")
		if err := stub.PutState(key, bs); err != nil {
			return fabric.SystemErrorResponse("Error putting key [" + key + "]: " + err.Error())
		}
		rows++
	}

	if p.Impuestos != nil {
		r, err := inscripciones.CommitPersonaImpuestos(stub, p.ID, p.Impuestos)
		if !err.IsOK() {
			return err
		}
		rows += r
	}

	return fabric.SuccessResponse("Ok", rows)
}
