package inscripciones

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func PutPersonaImpuestos(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {

	if len(args) != 2 {
		fabric.ClientErrorResponse("Numero incorrecto de parametros. Se espera 2 {CUIT, JSON}")
	}
	cuitStr := args[0]

	log.Print("args[0] [" + cuitStr + "]")
	if _, err := helpers.GetCUIT(args[0]); err != nil {
		return fabric.ClientErrorResponse("CUIT [" + cuitStr + "] invalida")
	}
	if exists, err := fabric.KeyExists(stub, cuitStr); !err.IsOK() {
		return err
	} else if !exists {
		return fabric.ClientErrorResponse("CUIT [" + cuitStr + "] inexistente")
	}

	impuestos := &model.Persona{}
	if err := json.Unmarshal([]byte(args[1]), &impuestos); err != nil {
		log.Print(args[1])
		return fabric.SystemErrorResponse("JSON invalido: " + err.Error())
	}

	rows, err := CommitPersonaImpuestos(stub, cuitStr, impuestos.Impuestos)
	if !err.IsOK() {
		log.Print(err.Msg)
		return err
	}
	log.Print(strconv.Itoa(rows) + " assets processed !!!")
	return fabric.SuccessResponse("Ok", rows)
}

func CommitPersonaImpuestos(stub shim.ChaincodeStubInterface, cuit string, kimps map[string]model.Impuesto) (int, *fabric.Response) {

	imps := []*model.Impuesto{}
	for _, imp := range kimps {
		imp := imp
		imps = append(imps, &imp)
	}

	if hid, impuestoDuplicado := impuestos.HasDuplicatedImpuestos(imps); hid {
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
		key := "PER_" + cuit + "_IMP_" + strconv.Itoa(int(imp.Impuesto))
		if err := stub.PutState(key, impuestoAsBytes); err != nil {
			return 0, fabric.SystemErrorResponse("Error putting key ["+key+"]: "+err.Error(), count)
		}
		count++
	}
	return len(imps), &fabric.Response{}
}
