package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) putPersona(APIstub shim.ChaincodeStubInterface, args []string, fType formatType) peer.Response {
	if len(args) != 2 {
		return shim.Error("Numero incorrecto de parametros. Se esperan 2 parametros {CUIT, JSON/PROTOBUF}")
	}
	isNotModeTest := !s.isModeTest
	var cuit uint64
	var err error
	cuitStr := args[0]
	if isNotModeTest {
		if err = checkClientID(APIstub); err != nil {
			return shim.Error(err.Error())
		}
	}

	log.Print("cuit recibido [" + cuitStr + "]")
	if cuit, err = getCUITArgs(args); err != nil {
		return shim.Error("CUIT [" + cuitStr + "] invalido")
	}
	var newPersona Persona
	if err = argToPersona([]byte(args[1]), &newPersona, fType); err != nil {
		return shim.Error("Formato invalido: " + err.Error())
	}
	return s.savePersona(APIstub, cuit, &newPersona)
}

func (s *SmartContract) putPersonas(APIstub shim.ChaincodeStubInterface, args []string, fType formatType) peer.Response {
	if len(args) != 1 {
		return shim.Error("Numero incorrecto de parametros. Se esperan 1 parametros {JSON/PROTOBUF}")
	}
	isNotModeTest := !s.isModeTest
	if isNotModeTest {
		if err := checkClientID(APIstub); err != nil {
			return shim.Error(err.Error())
		}
	}
	var newPersonas Personas
	if err := argToPersonas([]byte(args[0]), &newPersonas, fType); err != nil {
		return shim.Error("Formato invalido: " + err.Error())
	}

	rows := 0
	for _, p := range newPersonas.Personas {
		rows = rows + 1
		log.Printf("Grabando persona %d", p.CUIT)
		res := s.savePersona(APIstub, p.CUIT, p)
		if res.Status != shim.OK {
			return res
		}
	}
	msg := strconv.Itoa(rows) + " personas processed !!!"
	log.Print(msg)
	return shim.Success([]byte(msg))
}

func checkClientID(APIstub shim.ChaincodeStubInterface) error {
	// Get the client ID object
	id, err := cid.New(APIstub)
	if err != nil {
		return errors.New("Error at Get the client ID object [cid.New(APIstub)]")
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		return errors.New("Error at Get the client ID object [GetMSPID()]")
	}
	if mspid != "AFIP" {
		return errors.New("mspid [" + mspid + "] - La funcion putPersona solo puede ser invocada por AFIP")
	}
	return err
}

func (s *SmartContract) savePersona(APIstub shim.ChaincodeStubInterface, cuit uint64, newPersona *Persona) peer.Response {
	cuitStr := strconv.FormatUint(cuit, 10)
	if cuit != newPersona.CUIT {
		return shim.Error("El parametro cuit [" + cuitStr + "] y la cuit [" + strconv.FormatUint(newPersona.CUIT, 10) + "] en la Persona deben ser iguales")
	}
	switch newPersona.TipoPersona {
	case "F":
		if !strings.HasPrefix(cuitStr, "2") {
			return shim.Error("la cuit [" + cuitStr + "] debe comenzar con 2 cuando corresponde a una persona humana")
		}
		newPersona.Nombre = strings.ToUpper(strings.Trim(newPersona.Nombre, " "))
		newPersona.Apellido = strings.ToUpper(strings.Trim(newPersona.Apellido, " "))
		if len(newPersona.Apellido) == 0 {
			return shim.Error("persona humana sin apellido")
		}
		if len(newPersona.RazonSocial) > 0 {
			return shim.Error("persona humana con razonSocial [" + newPersona.RazonSocial + "]")
		}
		if newPersona.IDFormaJuridica != 0 {
			return shim.Error("persona humana con idFormaJuridica")
		}
		if !(newPersona.TipoDoc >= 1 && newPersona.TipoDoc <= 99) {
			return shim.Error("tipoDoc [" + strconv.Itoa(int(newPersona.TipoDoc)) + "] invalido, debe ser un entero entre 1 y 99")
		}
		newPersona.Documento = strings.ToUpper(strings.Trim(newPersona.Documento, " "))
		if len(newPersona.Documento) == 0 {
			return shim.Error("persona humana sin documento")
		}
		if len(newPersona.FechaNacimiento) == 0 {
			return shim.Error("persona humana sin fechaNacimiento")
		}
	case "J":
		if !strings.HasPrefix(cuitStr, "3") {
			return shim.Error("la cuit [" + cuitStr + "] debe comenzar con 3 cuando corresponde a una persona juridica")
		}
		newPersona.RazonSocial = strings.ToUpper(strings.Trim(newPersona.RazonSocial, " "))
		if len(newPersona.RazonSocial) == 0 {
			return shim.Error("persona juridica sin razonSocial")
		}
		if len(newPersona.Nombre) > 0 {
			return shim.Error("persona juridica con nombre [" + newPersona.Nombre + "], solo debe tener RazonSocial")
		}
		if len(newPersona.Apellido) > 0 {
			return shim.Error("persona juridica con apellido [" + newPersona.Apellido + "], solo debe tener RazonSocial")
		}
		if !(newPersona.IDFormaJuridica >= 1 && newPersona.IDFormaJuridica <= 999) {
			return shim.Error("idFormaJuridica [" + strconv.Itoa(int(newPersona.IDFormaJuridica)) + "] invalida")
		}
		if newPersona.TipoDoc != 0 {
			return shim.Error("persona juridica con tipoDoc")
		}
		if len(newPersona.Documento) > 0 {
			return shim.Error("persona juridica con documento")
		}
		if len(newPersona.FechaNacimiento) > 0 {
			return shim.Error("persona juridica con fechaNacimiento")
		}
		if len(newPersona.FechaFallecimiento) > 0 {
			return shim.Error("persona juridica con fechaFallecimiento")
		}
	default:
		return shim.Error("tipoPersona [" + newPersona.TipoPersona + "] invalido, debe ser F (Humana) o J (Juridica)")
	}
	switch newPersona.EstadoCUIT {
	case "A":
	case "I":
	default:
		return shim.Error("estadoCuit [" + newPersona.EstadoCUIT + "] invalido, debe ser A (Activa) o I (Inactiva)")
	}

	key := getPersonaKey(newPersona)

	if exist, err := keyExists(APIstub, cuitStr); err != nil {
		return s.businessErrorResponse(err.Error())
	} else if !exist {
		log.Print("Putting [" + cuitStr + "]...")
		if err := APIstub.PutState(cuitStr, []byte("{}")); err != nil {
			return shim.Error("Error putting cuitStr [" + cuitStr + "]: " + err.Error())
		}
	}

	var impuestos = newPersona.Impuestos
	newPersona.Impuestos = nil

	personaAsBytes, _ := json.Marshal(newPersona)

	log.Print("Putting [" + key + "]...")
	if err := APIstub.PutState(key, personaAsBytes); err != nil {
		return shim.Error("Error putting key [" + key + "]: " + err.Error())
	}

	if rows, err := putPersonaImpuestos(APIstub, cuitStr, impuestos); err != nil {
		log.Print(err.Error())
		return shim.Error(err.Error())
	} else {
		msg := strconv.Itoa(rows+1) + " assets processed !!!"
		log.Print(msg)
		return shim.Success([]byte(msg))
	}
}
