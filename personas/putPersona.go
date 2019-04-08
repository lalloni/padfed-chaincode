package personas

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/inscripciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func PutPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se esperan 2 (ID, PERSONA)")
	}
	var cuit uint64
	var err error
	var res *fabric.Response
	cuitStr := args[0]
	log.Print("cuit recibido [" + cuitStr + "]")
	if cuit, err = helpers.GetCUIT(cuitStr); err != nil {
		return fabric.ClientErrorResponse("ID [" + cuitStr + "] invalido")
	}
	newPersona := &model.Persona{}
	if res = ArgToPersona([]byte(args[1]), newPersona); !res.IsOK() {
		return res
	}
	return SavePersona(stub, cuit, newPersona)
}

func PutPersonas(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se espera 1 (PERSONAS)")
	}
	newPersonas := &model.Personas{}
	if err := ArgToPersonas([]byte(args[0]), newPersonas); !err.IsOK() {
		return err
	}

	rows := 0
	for _, p := range newPersonas.Personas {
		p := p
		log.Printf("Grabando persona %d", p.ID)
		res := SavePersona(stub, p.ID, &p)
		if res.Status != shim.OK {
			res.WrongItem = rows
			return res
		}
		rows++
	}
	log.Println(strconv.Itoa(rows) + " personas processed !!!")
	return fabric.SuccessResponse("Ok", rows)
}

func SavePersona(stub shim.ChaincodeStubInterface, cuit uint64, p *model.Persona) *fabric.Response {

	cuits := strconv.FormatUint(cuit, 10)
	if cuit != p.ID {
		return fabric.ClientErrorResponse("El parametro cuit [" + cuits + "] y la cuit [" + helpers.FormatCUIT(p.ID) + "] en la Persona deben ser iguales")
	}

	key := GetPersonaKey(p)

	exist, err := fabric.KeyExists(stub, cuits)
	if !err.IsOK() {
		return err
	}

	if !exist {
		if p.Tipo == "" {
			return fabric.ClientErrorResponse("No existe un asset [" + cuits + "] - Debe informarse los datos identificarios de la Persona")
		}
		log.Println("Putting [" + cuits + "]...")
		if err := stub.PutState(cuits, []byte("{}")); err != nil {
			return fabric.SystemErrorResponse("Error putting cuitStr [" + cuits + "]: " + err.Error())
		}
	}

	impuestos := p.Impuestos
	p.Impuestos = nil

	personaAsBytes, _ := json.Marshal(p)

	log.Println("Putting [" + key + "]...")
	if err := stub.PutState(key, personaAsBytes); err != nil {
		return fabric.SystemErrorResponse("Error putting key [" + key + "]: " + err.Error())
	}

	rows, err := inscripciones.CommitPersonaImpuestos(stub, cuits, impuestos)
	if !err.IsOK() {
		return err
	}

	if p.Tipo != "" {
		rows++
	}

	log.Print("Ok")

	return fabric.SuccessResponse("Ok", rows)
}
