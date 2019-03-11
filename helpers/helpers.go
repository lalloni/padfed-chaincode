package helpers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/lalloni/afip/cuit"
	"github.com/pkg/errors"
)

var dateRegexp = *regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)

func GetCUIT(s string) (uint64, error) {
	c, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Errorf("CUIT debe ser un número %q", s)
	}
	if !cuit.IsValid(uint64(c)) {
		return 0, errors.Errorf("CUIT inválido %q", s)
	}
	return uint64(c), nil
}

/*
hasDuplicatedImpuestos Chequea si en array impuestos existen impuestos con el mismo impuesto.
return
 - true si existen impuestos duplicados y el primer impuesto que se repite.
 - false si no existen impuestos duplicados.

*/

func GetTxConfirmableKey(idOrganismo int, idTxc uint64) string {
	idOrganismoStr := strconv.Itoa(idOrganismo)
	idTxcStr := strconv.FormatUint(idTxc, 10)
	return "ORG_" + idOrganismoStr + "_TXC_" + idTxcStr
}

// func findPersona(stub shim.ChaincodeStubInterface, cuit uint64) (Persona, []byte, error) {
// 	var cuitStr = strconv.FormatUint(cuit, 10)
// 	personaAsBytes, err := stub.GetState("PER_" + cuitStr)
// 	var persona Persona
// 	if err != nil {
// 		return persona, personaAsBytes, errors.New("Error al buscar la Persona " + cuitStr)
// 	} else if personaAsBytes == nil {
// 		return persona, personaAsBytes, errors.New("No existe Persona para " + cuitStr)
// 	}
// 	err = argToPersona(personaAsBytes, &persona, JSON)
// 	return persona, personaAsBytes, nil
// }

func ValidateDate(dateStr string) error {
	if dateStr != "" {
		res := dateRegexp.FindStringSubmatch(dateStr)
		if len(res) != 4 {
			return fmt.Errorf("fecha inválida %s", dateStr)
		}
		y, _ := strconv.Atoi(res[1])
		m, _ := strconv.Atoi(res[2])
		d, _ := strconv.Atoi(res[3])
		date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		if date.Year() != y || int(date.Month()) != m || date.Day() != d {
			return fmt.Errorf("fecha inválida, ingresada [%s-%s-%s], convertida a time [%d-%d-%d]", res[1], res[2], res[3], date.Year(), int(date.Month()), date.Day())
		}
	}
	return nil
}

// TODO: eliminar generación manual de JSON
func WriteInBuffer(buffer *bytes.Buffer, bufferValue []byte, keyParameter string, bArrayMemberAlreadyWritten bool) {
	// Add a comma before array members, suppress it for the first array member
	if bArrayMemberAlreadyWritten {
		(*buffer).WriteString(",")
	}
	(*buffer).WriteString("{\"Key\":")
	(*buffer).WriteString("\"")
	(*buffer).WriteString(keyParameter)
	(*buffer).WriteString("\"")
	(*buffer).WriteString(",\"Record\":")
	(*buffer).Write(bufferValue)
	(*buffer).WriteString("}")
}
