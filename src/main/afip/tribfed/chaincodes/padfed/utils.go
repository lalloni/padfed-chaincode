package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	cuitVerifier "github.com/lalloni/afip/cuit"
	"github.com/xeipuuv/gojsonschema"
)

var (
	personaSchemaLoader  = gojsonschema.NewStringLoader(PersonaSchema)
	personasSchemaLoader = gojsonschema.NewStringLoader(PersonasSchema)
	impuestoSchemaLoader = gojsonschema.NewStringLoader(ImpuestoSchema)
)

var CHECK_DATE_REGEXP = *regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)

func getCUITArgs(args []string) (uint64, error) {
	var cuit int64
	var err error
	if cuit, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return 0, errors.New("CUIT debe ser un número " + args[0])
	}

	if !cuitVerifier.IsValid(uint64(cuit)) {
		return 0, errors.New("CUIT invalido " + args[0])
	}
	return uint64(cuit), nil
}

func argToPersona(personaAsBytes []byte, persona *Persona, fType formatType) Response {
	//fmt.Println(PersonaSchema)
	switch fType {
	case JSON:
		documentLoader := gojsonschema.NewStringLoader(string(personaAsBytes))
		result, err := gojsonschema.Validate(personaSchemaLoader, documentLoader)
		if err != nil {
			return clientErrorResponse("JSON schema invalido: " + err.Error() + " - " + string(personaAsBytes))
		}

		if !result.Valid() {
			var errosStr string
			for _, desc := range result.Errors() {
				errosStr += desc.Description() + ". "
			}
			return clientErrorResponse("JSON no cumple con el esquema: " + errosStr)
		}
		err = json.Unmarshal(personaAsBytes, persona)
		if err != nil {
			return systemErrorResponse("JSON invalido: " + err.Error())
		}
	case PROTOBUF:
		err := proto.Unmarshal(personaAsBytes, persona)
		if err != nil {
			log.Fatal("PROTOBUF invalido: ", err)
		}
	}
	return validatePersona(persona)
}

func argToPersonas(personasAsBytes []byte, personas *Personas, fType formatType) Response {
	switch fType {
	case JSON:
		documentLoader := gojsonschema.NewStringLoader(string(personasAsBytes))
		result, err := gojsonschema.Validate(personasSchemaLoader, documentLoader)
		if err != nil {
			return clientErrorResponse("JSON schema invalido: " + err.Error() + " - " + string(personasAsBytes))
		}
		if !result.Valid() {
			var errosStr string
			for _, desc := range result.Errors() {
				errosStr += desc.Description() + ". "
			}
			return clientErrorResponse("JSON no cumple con el esquema: " + errosStr)
		}
		err = json.Unmarshal(personasAsBytes, &personas)
		if err != nil {
			return systemErrorResponse("JSON invalido: " + err.Error())
		}
	case PROTOBUF:
		err := proto.Unmarshal(personasAsBytes, personas)
		if err != nil {
			log.Fatal("PROTOBUF invalido: ", err)
		}
	}

	for _, p := range personas.GetPersonas() {
		err := validatePersona(p)
		if err.isError() {
			return err
		}
	}
	return Response{}
}

func validatePersona(persona *Persona) Response {
	var err Response
	cuitStr := strconv.FormatUint(persona.CUIT, 10)
	if !cuitVerifier.IsValid(persona.CUIT) {
		return clientErrorResponse("cuit [" + cuitStr + "] invalida")
	}
	if err := validateDate(persona.Nacimiento); err != nil {
		return clientErrorResponse("nacimiento [" + persona.Nacimiento + "] invalido: " + err.Error())
	}
	if err := validateDate(persona.Inscripcion); err != nil {
		return clientErrorResponse("inscripcion [" + persona.Inscripcion + "] invalida: " + err.Error())
	}
	if err := validateDate(persona.FechaCierre); err != nil {
		return clientErrorResponse("fechaCierre [" + persona.FechaCierre + "] invalida: " + err.Error())
	}
	if err := validateDate(persona.Fallecimiento); err != nil {
		return clientErrorResponse("fallecimiento [" + persona.Fallecimiento + "] invalido: " + err.Error())
	}
	if err := validateDate(persona.DS); err != nil {
		return clientErrorResponse("ds [" + persona.DS + "] invalido: " + err.Error())
	}
	return err
}

/*
hasDuplicatedImpuestos Chequea si en array impuestos existen impuestos con el mismo impuesto.
return
 - true si existen impuestos duplicados y el primer impuesto que se repite.
 - false si no existen impuestos duplicados.

*/
func hasDuplicatedImpuestos(impuestos []*Impuesto) (bool, *Impuesto) {
	var index = make(map[int]*Impuesto)
	for _, imp := range impuestos {
		if _, exist := index[int(imp.Impuesto)]; exist {
			return exist, imp
		} else {
			index[int(imp.Impuesto)] = imp
		}
	}
	return false, &Impuesto{}
}

func getPersonaKey(persona *Persona) string {
	cuitStr := strconv.FormatUint(persona.CUIT, 10)
	return "PER_" + cuitStr
}

func getImpuestoKeyByCuitId(cuit uint64, impuesto int32) string {
	cuitStr := strconv.FormatUint(cuit, 10)
	impStr := strconv.Itoa(int(impuesto))
	return "PER_" + cuitStr + "_IMP_" + impStr
}

func getTxConfirmableKey(idOrganismo int, idTxc uint64) string {
	idOrganismoStr := strconv.Itoa(idOrganismo)
	idTxcStr := strconv.FormatUint(idTxc, 10)
	return "ORG_" + idOrganismoStr + "_TXC_" + idTxcStr
}

// func findPersona(APIstub shim.ChaincodeStubInterface, cuit uint64) (Persona, []byte, error) {
// 	var cuitStr = strconv.FormatUint(cuit, 10)
// 	personaAsBytes, err := APIstub.GetState("PER_" + cuitStr)
// 	var persona Persona
// 	if err != nil {
// 		return persona, personaAsBytes, errors.New("Error al buscar la Persona " + cuitStr)
// 	} else if personaAsBytes == nil {
// 		return persona, personaAsBytes, errors.New("No existe Persona para " + cuitStr)
// 	}
// 	err = argToPersona(personaAsBytes, &persona, JSON)
// 	return persona, personaAsBytes, nil
// }

func findImpuesto(APIstub shim.ChaincodeStubInterface, cuit uint64, idImpuesto int32) (Impuesto, []byte, error) {
	var cuitStr = strconv.FormatUint(cuit, 10)
	var impuestoStr = strconv.Itoa(int(idImpuesto))

	impuestoAsBytes, err := APIstub.GetState(getImpuestoKeyByCuitId(cuit, idImpuesto))
	var impuesto Impuesto
	if err != nil {
		return impuesto, impuestoAsBytes, errors.New("Error al buscar Impuesto " + cuitStr)
	} else if impuestoAsBytes == nil {
		return impuesto, impuestoAsBytes, errors.New("No existe Impuesto para la CUIT " + cuitStr + " e impuesto " + impuestoStr)
	}
	err = json.Unmarshal(impuestoAsBytes, &impuesto)
	return impuesto, impuestoAsBytes, nil
}

func validateDate(dateStr string) error {
	if dateStr != "" {
		res := CHECK_DATE_REGEXP.FindStringSubmatch(dateStr)
		if len(res) != 4 {
			return fmt.Errorf("Fecha invalida %s", dateStr)
		}
		y, _ := strconv.Atoi(res[1])
		m, _ := strconv.Atoi(res[2])
		d, _ := strconv.Atoi(res[3])
		date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

		if date.Year() != y || int(date.Month()) != m || date.Day() != d {
			return fmt.Errorf("Fecha invalida, ingresada [%s-%s-%s], convertida a time [%d-%d-%d]", res[1], res[2], res[3], date.Year(), int(date.Month()), date.Day())
		}
	}
	return nil
}

// keyExists returns true if the key exists
func keyExists(APIstub shim.ChaincodeStubInterface, key string) (bool, Response) {
	log.Println("Key[" + key + "] using GetState...")
	if assetAsByte, err := APIstub.GetState(key); err != nil {
		return false, systemErrorResponse(err.Error())
	} else {
		return assetAsByte != nil, Response{}
	}
}

func writeInBuffer(buffer *bytes.Buffer, bufferValue string, keyParameter string, bArrayMemberAlreadyWritten bool) {
	// Add a comma before array members, suppress it for the first array member
	if bArrayMemberAlreadyWritten == true {
		(*buffer).WriteString(",")
	}
	(*buffer).WriteString("{\"Key\":")
	(*buffer).WriteString("\"")
	(*buffer).WriteString(keyParameter)
	(*buffer).WriteString("\"")
	(*buffer).WriteString(",\"Record\":")
	(*buffer).WriteString(bufferValue)
	(*buffer).WriteString("}")
}
