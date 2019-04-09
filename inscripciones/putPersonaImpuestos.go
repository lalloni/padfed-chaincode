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

func CommitPersonaImpuestos(stub shim.ChaincodeStubInterface, cuit uint64, kimps map[string]*model.PersonaImpuesto) (int, *fabric.Response) {

	count := 0
	set := map[uint]struct{}{}
	for _, v := range kimps {
		if v == nil {
			continue
		}
		if _, ok := set[v.Impuesto]; ok {
			return 0, fabric.ClientErrorResponse(fmt.Sprintf("Impuesto %d duplicado", v.Impuesto), count)
		}
		count++
	}

	count = 0
	for key, imp := range kimps {
		if imp == nil {
			// del
			impuesto, err := strconv.ParseUint(key, 10, strconv.IntSize)
			if err != nil {
				return count, fabric.ClientErrorResponse(fmt.Sprintf("Código de impuesto inválido: %q", key), count)
			}
			res := DeletePersonaImpuesto(stub, cuit, uint(impuesto))
			if !res.IsOK() {
				res.WrongItem = count
				return count, res
			}
		} else {
			// put
			res := PutPersonaImpuesto(stub, cuit, imp)
			if !res.IsOK() {
				res.WrongItem = count
				return count, res
			}
		}
		count++
	}

	return count, &fabric.Response{}
}

func DeletePersonaImpuesto(stub shim.ChaincodeStubInterface, cuit uint64, impuesto uint) *fabric.Response {
	exist, errr := impuestos.ExistsIDImpuesto(stub, impuesto)
	if !errr.IsOK() {
		return errr
	}
	if !exist {
		return fabric.ClientErrorResponse(fmt.Sprintf("No existe el impuesto %d en la persona %d", impuesto, cuit))
	}
	err := stub.DelState(GetImpuestoKeyByCuitID(cuit, impuesto))
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("Eliminando impuesto de persona: %v", err))
	}
	return &fabric.Response{}
}

func PutPersonaImpuesto(stub shim.ChaincodeStubInterface, cuit uint64, imp *model.PersonaImpuesto) *fabric.Response {
	exists, res := impuestos.ExistsIDImpuesto(stub, imp.Impuesto)
	if !res.IsOK() {
		return res
	}
	if !exists {
		return fabric.ClientErrorResponse("Impuesto [" + strconv.Itoa(int(imp.Impuesto)) + "] inexistente")
	}
	bs, err := json.Marshal(imp)
	if err != nil {
		return fabric.SystemErrorResponse("Error marshalling impuesto [" + strconv.Itoa(int(imp.Impuesto)) + "]: " + err.Error())
	}
	key := GetImpuestoKeyByCuitID(cuit, imp.Impuesto)
	if err := stub.PutState(key, bs); err != nil {
		return fabric.SystemErrorResponse("Error putting key [" + key + "]: " + err.Error())
	}
	return &fabric.Response{}
}
