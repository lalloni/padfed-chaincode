package personas

import (
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
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
	newPersonas := []model.Persona{}
	if err := model.ArgToPersonas([]byte(args[0]), &newPersonas); !err.IsOK() {
		return err
	}

	rows := 0
	for _, p := range newPersonas {
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
