package impuestos

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func PutParamImpuestos(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de parámetros. Se espera un parámetro con {[]JSON}")
	}
	var err error

	var impuestos []model.Impuesto
	if err = json.Unmarshal([]byte(args[0]), &impuestos); err != nil {
		msg := "JSON invalido: " + err.Error()
		log.Println(msg)
		return fabric.SystemErrorResponse(msg)
	}
	count := 0
	for _, imp := range impuestos {
		if !(imp.IDOrg == 1 || (imp.IDOrg >= 900 && imp.IDOrg <= 999)) {
			return fabric.ClientErrorResponse("idOrg ["+strconv.Itoa(int(imp.IDOrg))+"] debe ser un entero igual a 1:AFIP o entre 900 y 999", count)
		}

		if response := ValidateIDImpuesto(imp.Impuesto); !response.IsOK() {
			response.WrongItem = count
			return response
		}

		impuestoAsBytes, _ := json.Marshal(imp)
		var key = GetParamImpuestoKey(imp.Impuesto)
		if err := stub.PutState(key, impuestoAsBytes); err != nil {
			return fabric.SystemErrorResponse("Error putting key ["+key+"] - "+err.Error(), count)
		}
		count++
	}
	log.Println(strconv.Itoa(len(impuestos)) + " assets processed !!!")
	return fabric.SuccessResponse("Ok", len(impuestos))
}

// ValidateIDImpuesto valida el codigo de impuesto sin acceder al State
func ValidateIDImpuesto(impuesto uint) *fabric.Response {
	if !(impuesto >= 1 && impuesto <= 9999) {
		return fabric.ClientErrorResponse("impuesto [" + strconv.Itoa(int(impuesto)) + "] debe ser un entero entre 1 y 9999")
	}
	return &fabric.Response{}
}

// ExistsIDImpuesto verifica que exista un asset "IMP_<idImpuesot>"
func ExistsIDImpuesto(stub shim.ChaincodeStubInterface, impuesto uint) (bool, *fabric.Response) {
	exists, response := fabric.KeyExists(stub, GetParamImpuestoKey(impuesto))
	if !response.IsOK() {
		return false, response
	}
	return exists, &fabric.Response{}
}

func GetParamImpuestoKey(impuesto uint) string {
	return "IMP_" + strconv.Itoa(int(impuesto))
}
