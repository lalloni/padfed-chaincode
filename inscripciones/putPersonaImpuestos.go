package inscripciones

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func PutPersonaImpuestos(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {

	if len(args) != 1 {
		fabric.ClientErrorResponse("Numero incorrecto de parametros. Se espera 1.")
	}
	p := &model.Persona{}
	err := model.ArgToPersona([]byte(args[0]), p)
	if !err.IsOK() {
		return err
	}
	if exists, err := fabric.KeyExists(stub, strconv.Itoa(int(p.ID))); !err.IsOK() {
		return err
	} else if !exists {
		return fabric.ClientErrorResponse(fmt.Sprintf("CUIT [%d] inexistente", p.ID))
	}

	rows, err := CommitPersonaImpuestos(stub, p.ID, p.Impuestos)
	if !err.IsOK() {
		return err
	}
	log.Print(strconv.Itoa(rows) + " assets processed !!!")
	return fabric.SuccessResponse("Ok", rows)
}

func CommitPersonaImpuestos(stub shim.ChaincodeStubInterface, cuit uint64, kimps map[string]model.PersonaImpuesto) (int, *fabric.Response) {

	imps := []*model.PersonaImpuesto{}
	for _, imp := range kimps {
		imp := imp
		imps = append(imps, &imp)
	}

	if hid, impuestoDuplicado := HasDuplicatedImpuestos(imps); hid {
		return 0, fabric.ClientErrorResponse("Array con impuesto [" + strconv.Itoa(int(impuestoDuplicado.Impuesto)) + "] duplicado")
	}
	count := 0
	for _, imp := range imps {
		if exists, err := impuestos.ExistsIDImpuesto(stub, imp.Impuesto); !err.IsOK() {
			err.WrongItem = count
			return 0, err
		} else if !exists {
			return 0, fabric.ClientErrorResponse("impuesto ["+strconv.Itoa(int(imp.Impuesto))+"] no definido en ParamImpuesto", count)
		}
		impuestoAsBytes, err := json.Marshal(imp)
		if err != nil {
			return 0, fabric.SystemErrorResponse("Error marshalling impuesto ["+strconv.Itoa(int(imp.Impuesto))+"]: "+err.Error(), count)
		}
		key := GetImpuestoKeyByCuitID(cuit, imp.Impuesto)
		if err := stub.PutState(key, impuestoAsBytes); err != nil {
			return 0, fabric.SystemErrorResponse("Error putting key ["+key+"]: "+err.Error(), count)
		}
		count++
	}
	return len(imps), &fabric.Response{}
}
