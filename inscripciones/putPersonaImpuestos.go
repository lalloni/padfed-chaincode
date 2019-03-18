package inscripciones

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

var periodoFiscalRegexp = *regexp.MustCompile(`^(19\d{2}|20[0-2]\d|2030)(0\d|1[0-2])$`)

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

	impuestos := &model.Impuestos{}
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

func CommitPersonaImpuestos(stub shim.ChaincodeStubInterface, cuit string, imps []*model.Impuesto) (int, *fabric.Response) {

	if hid, impuestoDuplicado := impuestos.HasDuplicatedImpuestos(imps); hid {
		return 0, fabric.ClientErrorResponse("Array con impuesto [" + strconv.Itoa(int(impuestoDuplicado.Impuesto)) + "] duplicado")
	}
	count := 0
	for _, imp := range imps {
		impuestoAsBytes, _ := json.Marshal(imp)

		if err := impuestos.ValidateIDImpuesto(imp.Impuesto); !err.IsOK() {
			err.WrongItem = count
			return 0, err
		}
		if exists, err := impuestos.ExistsIDImpuesto(stub, imp.Impuesto); !err.IsOK() {
			err.WrongItem = count
			return 0, err
		} else if !exists {
			return 0, fabric.ClientErrorResponse("impuesto ["+strconv.Itoa(int(imp.Impuesto))+"] no definido en ParamImpuesto", count)
		}
		if err := helpers.ValidateDate(imp.Inscripcion); err != nil {
			return 0, fabric.ClientErrorResponse("inscripcion ["+imp.Inscripcion+"]: "+err.Error(), count)
		}
		if err := helpers.ValidateDate(imp.DS); err != nil {
			return 0, fabric.ClientErrorResponse("ds ["+imp.DS+"]: "+err.Error(), count)
		}
		periodoString := strconv.FormatInt(int64(imp.Periodo), 10)
		res := periodoFiscalRegexp.FindStringSubmatch(periodoString)
		if len(res) != 3 {
			return 0, fabric.ClientErrorResponse("periodo ["+strconv.Itoa(int(imp.Periodo))+"] debe tener formato YYYY00 o YYYYMM con YYYY entre 1900 y 2030", count)
		}
		if (imp.Dia < 0) || (imp.Dia > 31) {
			return 0, fabric.ClientErrorResponse("dia ["+strconv.Itoa(int(imp.Dia))+"] debe ser un entero entre 1 y 31 o nulo", count)
		}
		key := "PER_" + cuit + "_IMP_" + strconv.Itoa(int(imp.Impuesto))
		if err := stub.PutState(key, impuestoAsBytes); err != nil {
			return 0, fabric.SystemErrorResponse("Error putting key ["+key+"]: "+err.Error(), count)
		}
		switch imp.Estado {
		case "AC", "AT", "BP", "BD", "EX", "ET", "NA":
		default:
			return 0, fabric.ClientErrorResponse("estado ["+imp.Estado+"] invalido, debe ser AC (Activo), AT (Activo en tramite), BP (Baja provisoria), BD (Baja definitiva), EX (Exento), ET (Exento en tramite), NA (No aportante)", count)
		}
		count++
	}
	return len(imps), &fabric.Response{}
}
