package personas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func GetPersonaAPI(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	cuit, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fabric.ClientErrorResponse(fmt.Sprintf("cuit inválido: %v", args[0]))
	}
	p, err := GetPersona(stub, cuit)
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("obteniendo persona: %v", err))
	}
	if p == nil {
		return fabric.NotFoundErrorResponse()
	}
	bs, err := json.Marshal(p)
	if err != nil {
		return fabric.SystemErrorResponse(fmt.Sprintf("generando respuesta: %v", err))
	}
	return fabric.SuccessResponseWithBuffer(bytes.NewBuffer(bs))
}

func GetPersona(stub shim.ChaincodeStubInterface, cuit uint64) (*model.Persona, error) {

	p := &model.Persona{}

	key := GetPersonaKeyCUIT(cuit)

	it, err := stub.GetStateByRange(key, key+"z")
	if err != nil {
		return nil, errors.Wrap(err, "getting persona values")
	}
	p.ID = cuit
	defer it.Close()

	if !it.HasNext() {
		return nil, nil
	}

	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, errors.Wrap(err, "reading next key-value")
		}
		switch {
		case key == kv.GetKey():
			b := &model.PersonaBasica{}
			err := json.Unmarshal(kv.GetValue(), b)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshaling persona")
			}
			p.Persona = b
		default: // assume que es impuesto
			i := &model.PersonaImpuesto{}
			err := json.Unmarshal(kv.GetValue(), i)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshaling persona impuesto")
			}
			if p.Impuestos == nil {
				p.Impuestos = map[string]*model.PersonaImpuesto{}
			}
			p.Impuestos[strconv.Itoa(int(i.Impuesto))] = i
		}
	}

	return p, nil
}
