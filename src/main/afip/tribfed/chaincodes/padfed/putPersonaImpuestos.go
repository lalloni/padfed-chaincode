package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var CHECK_PERIODO_FISCAL_REGEXP = *regexp.MustCompile(`^(19\d{2}|20[0-2]\d|2030)(0\d|1[0-2])$`)

func (s *SmartContract) putPersonaImpuestos(APIstub shim.ChaincodeStubInterface, args []string) Response {

	if len(args) != 2 {
		clientErrorResponse("Numero incorrecto de parametros. Se espera 2 {CUIT, JSON}")
	}
	cuitStr := args[0]

	log.Print("args[0] [" + cuitStr + "]")
	if _, err := getCUITArgs(args); err != nil {
		return clientErrorResponse("CUIT [" + cuitStr + "] invalida")
	}
	if exists, err := keyExists(APIstub, cuitStr); err.isError() {
		return err
	} else if !exists {
		return clientErrorResponse("CUIT [" + cuitStr + "] inexistente")
	}

	var impuestos Impuestos
	if err := json.Unmarshal([]byte(args[1]), &impuestos); err != nil {
		log.Print(args[1])
		return systemErrorResponse("JSON invalido: " + err.Error())
	}

	if rows, err := s.commitPersonaImpuestos(APIstub, cuitStr, impuestos.Impuestos); err.isError() {
		log.Print(err.Msg)
		return err
	} else {
		log.Print(strconv.Itoa(rows) + " assets processed !!!")
		return successResponse("Ok", rows)
	}
}

func (s *SmartContract) commitPersonaImpuestos(APIstub shim.ChaincodeStubInterface, cuit string, impuestos []*Impuesto) (int, Response) {

	if hid, impuestoDuplicado := hasDuplicatedImpuestos(impuestos); hid {
		return 0, clientErrorResponse("Array con impuesto [" + strconv.Itoa(int(impuestoDuplicado.Impuesto)) + "] duplicado")
	}
	count := 0
	for _, imp := range impuestos {
		impuestoAsBytes, _ := json.Marshal(imp)

		if !(imp.IDOrganismo == 0 || imp.IDOrganismo == 1 || (imp.IDOrganismo >= 900 && imp.IDOrganismo <= 999)) {
			return 0, clientErrorResponse("idOrg ["+strconv.Itoa(int(imp.IDOrganismo))+"] must be an integer 1:AFIP or between 900 and 999", count)
		}
		if err := validateIdImpuesto(imp.Impuesto); err.isError() {
			err.WrongItem = count
			return 0, err
		}
		if exists, err := existsIdImpuesto(APIstub, imp.Impuesto); err.isError() {
			err.WrongItem = count
			return 0, err
		} else if !exists {
			return 0, clientErrorResponse("impuesto ["+strconv.Itoa(int(imp.Impuesto))+"] no definido en ParamImpuesto", count)
		}
		if err := validateDate(imp.Inscripcion); err != nil {
			return 0, clientErrorResponse("inscripcion ["+imp.Inscripcion+"]: "+err.Error(), count)
		}
		if err := validateDate(imp.DS); err != nil {
			return 0, clientErrorResponse("ds ["+imp.DS+"]: "+err.Error(), count)
		}
		periodoString := strconv.FormatInt(int64(imp.Periodo), 10)
		res := CHECK_PERIODO_FISCAL_REGEXP.FindStringSubmatch(periodoString)
		if len(res) != 3 {
			return 0, clientErrorResponse("periodo ["+strconv.Itoa(int(imp.Periodo))+"] debe tener formato YYYY00 o YYYYMM con YYYY entre 1900 y 2030", count)
		}
		if (imp.Dia < 0) || (imp.Dia > 31) {
			return 0, clientErrorResponse("dia ["+strconv.Itoa(int(imp.Dia))+"] debe ser un entero entre 1 y 31 o nulo", count)
		}
		key := "PER_" + cuit + "_IMP_" + strconv.Itoa(int(imp.Impuesto))
		if err := APIstub.PutState(key, impuestoAsBytes); err != nil {
			return 0, systemErrorResponse("Error putting key ["+key+"]: "+err.Error(), count)
		}
		switch imp.Estado {
		case "AC", "AT", "BP", "BD", "EX", "ET", "NA":
		default:
			return 0, clientErrorResponse("estado ["+imp.Estado+"] invalido, debe ser AC (Activo), AT (Activo en tramite), BP (Baja provisoria), BD (Baja definitiva), EX (Exento), ET (Exento en tramite), NA (No aportante)", count)
		}
		count++
	}
	return len(impuestos), Response{}
}
