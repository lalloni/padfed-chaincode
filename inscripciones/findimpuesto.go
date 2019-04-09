package inscripciones

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func FindImpuesto(stub shim.ChaincodeStubInterface, cuit uint64, idImpuesto uint) (model.Impuesto, []byte, error) {
	var cuitStr = strconv.FormatUint(cuit, 10)
	var impuestoStr = strconv.Itoa(int(idImpuesto))

	impuestoAsBytes, err := stub.GetState(GetImpuestoKeyByCuitID(cuit, idImpuesto))
	var impuesto model.Impuesto
	if err != nil {
		return impuesto, impuestoAsBytes, errors.New("Error al buscar Impuesto " + cuitStr)
	} else if impuestoAsBytes == nil {
		return impuesto, impuestoAsBytes, errors.New("No existe Impuesto para la CUIT " + cuitStr + " e impuesto " + impuestoStr)
	}
	err = json.Unmarshal(impuestoAsBytes, &impuesto)
	return impuesto, impuestoAsBytes, err
}
