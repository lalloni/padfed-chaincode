package personas

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
)

func DelPersonaAssets(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Numero incorrecto de parametros. Se esperaba 2 parametros con {CUIL, []KEYS}")
	}
	var cuit = args[0]
	var keys []string
	if err := json.Unmarshal([]byte(args[1]), &keys); err != nil {
		msg := "JSON de array invalido: " + err.Error()
		log.Println(msg)
		return fabric.SystemErrorResponse(msg)
	}

	if len(keys) == 0 {
		return fabric.ClientErrorResponse("El array debe tener por lo menos un elemento")
	}
	if errResponse := checkDuplicated(keys, cuit); errResponse != nil {
		return errResponse
	}
	exists, err := fabric.KeyExists(stub, "PER_"+cuit)
	if !err.IsOK() {
		return err
	}
	if !exists {
		return fabric.ClientErrorResponse("Debe existir como un asset la persona: " + cuit)
	}
	count := 0
	for _, key := range keys {
		fabric.Log.Infof("key to delete [%s]", key)
		if err := stub.DelState(key); err != nil {
			return fabric.SystemErrorResponse("Error al eliminar: [" + key + "] " + err.Error())
		}
		count++
	}

	fabric.Log.Infof("%d assets deleted!!!", count)
	return fabric.SuccessResponse("Ok", count)
}

var re = regexp.MustCompile(`^(PER_)(\d{11})(_IMP_)(\d+)$`)

func checkDuplicated(array []string, cuit string) *fabric.Response {
	keys := make(map[string]bool)
	for _, entry := range array {
		if _, ok := keys[entry]; ok {
			return fabric.ClientErrorResponse("El array no puede tener elementos repetidos")
		}
		keys[entry] = true
		res := re.FindStringSubmatch(entry)
		if res == nil || res[2] != cuit {
			return fabric.ClientErrorResponse("El array debe tener keys compuestas de assets PersonaImpuesto PER_" + cuit + "_IMP_[imp]")
		}
	}
	return nil
}
