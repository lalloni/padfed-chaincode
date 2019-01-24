package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s *SmartContract) putPersona(APIstub shim.ChaincodeStubInterface, args []string, fType formatType) Response {
	if len(args) != 2 {
		return clientErrorResponse("Numero incorrecto de parametros. Se esperan 2 parametros {CUIT, JSON/PROTOBUF}")
	}
	isNotModeTest := !s.isModeTest
	var cuit uint64
	var err error
	var err_R Response
	cuitStr := args[0]
	if isNotModeTest {
		if err_R = s.checkClientID(); err_R == (Response{}) {
			return err_R
		}
	}
	log.Print("cuit recibido [" + cuitStr + "]")
	if cuit, err = getCUITArgs(args); err != nil {
		return clientErrorResponse("CUIT [" + cuitStr + "] invalido")
	}
	var newPersona Persona
	if err_R = argToPersona([]byte(args[1]), &newPersona, fType); err_R == (Response{}) {
		return err_R
	}
	return s.savePersona(APIstub, cuit, &newPersona)
}

func (s *SmartContract) putPersonas(APIstub shim.ChaincodeStubInterface, args []string, fType formatType) Response {
	if len(args) != 1 {
		return clientErrorResponse("Numero incorrecto de parametros. Se esperan 1 parametros {JSON/PROTOBUF}")
	}
	isNotModeTest := !s.isModeTest
	if isNotModeTest {
		if err := s.checkClientID(); err == (Response{}) {
			return err
		}
	}
	var newPersonas Personas
	if err := argToPersonas([]byte(args[0]), &newPersonas, fType); err == (Response{}) {
		return err
	}

	rows := 0
	for _, p := range newPersonas.Personas {
		log.Printf("Grabando persona %d", p.CUIT)
		res := s.savePersona(APIstub, p.CUIT, p)
		if res.Status != shim.OK {
			res.WrongItem = rows
			return res
		}
		rows++
	}
	log.Println(strconv.Itoa(rows) + " personas processed !!!")
	return successResponse("Ok", rows)
}

func (s *SmartContract) checkClientID() Response {
	if s.mspid != "AFIP" {
		return forbiddenErrorResponse("mspid [" + s.mspid + "] - La funcion putPersona solo puede ser invocada por AFIP")
	}
	return Response{}
}

func (s *SmartContract) savePersona(APIstub shim.ChaincodeStubInterface, cuit uint64, newPersona *Persona) Response {
	tipoPersonaNull := false
	cuitStr := strconv.FormatUint(cuit, 10)
	if cuit != newPersona.CUIT {
		return clientErrorResponse("El parametro cuit [" + cuitStr + "] y la cuit [" + strconv.FormatUint(newPersona.CUIT, 10) + "] en la Persona deben ser iguales")
	}
	switch newPersona.TipoPersona {
	case "F":
		if !strings.HasPrefix(cuitStr, "2") {
			return clientErrorResponse("la cuit [" + cuitStr + "] debe comenzar con 2 cuando corresponde a una persona humana")
		}
		newPersona.Nombre = strings.ToUpper(strings.Trim(newPersona.Nombre, " "))
		newPersona.Apellido = strings.ToUpper(strings.Trim(newPersona.Apellido, " "))
		if len(newPersona.Apellido) == 0 {
			return clientErrorResponse("persona humana sin apellido")
		}
		if len(newPersona.RazonSocial) > 0 {
			return clientErrorResponse("persona humana con razonSocial [" + newPersona.RazonSocial + "]")
		}
		if newPersona.IDFormaJuridica != 0 {
			return clientErrorResponse("persona humana con idFormaJuridica")
		}
		if !(newPersona.TipoDoc >= 1 && newPersona.TipoDoc <= 99) {
			return clientErrorResponse("tipoDoc [" + strconv.Itoa(int(newPersona.TipoDoc)) + "] invalido, debe ser un entero entre 1 y 99")
		}
		newPersona.Documento = strings.ToUpper(strings.Trim(newPersona.Documento, " "))
		if len(newPersona.Documento) == 0 {
			return clientErrorResponse("persona humana sin documento")
		}
		if len(newPersona.FechaNacimiento) == 0 {
			return clientErrorResponse("persona humana sin fechaNacimiento")
		}
	case "J":
		if !strings.HasPrefix(cuitStr, "3") {
			return clientErrorResponse("la cuit [" + cuitStr + "] debe comenzar con 3 cuando corresponde a una persona juridica")
		}
		newPersona.RazonSocial = strings.ToUpper(strings.Trim(newPersona.RazonSocial, " "))
		if len(newPersona.RazonSocial) == 0 {
			return clientErrorResponse("persona juridica sin razonSocial")
		}
		if len(newPersona.Nombre) > 0 {
			return clientErrorResponse("persona juridica con nombre [" + newPersona.Nombre + "], solo debe tener RazonSocial")
		}
		if len(newPersona.Apellido) > 0 {
			return clientErrorResponse("persona juridica con apellido [" + newPersona.Apellido + "], solo debe tener RazonSocial")
		}
		if !(newPersona.IDFormaJuridica >= 1 && newPersona.IDFormaJuridica <= 999) {
			return clientErrorResponse("idFormaJuridica [" + strconv.Itoa(int(newPersona.IDFormaJuridica)) + "] invalida")
		}
		if newPersona.TipoDoc != 0 {
			return clientErrorResponse("persona juridica con tipoDoc")
		}
		if len(newPersona.Documento) > 0 {
			return clientErrorResponse("persona juridica con documento")
		}
		if len(newPersona.FechaNacimiento) > 0 {
			return clientErrorResponse("persona juridica con fechaNacimiento")
		}
		if len(newPersona.FechaFallecimiento) > 0 {
			return clientErrorResponse("persona juridica con fechaFallecimiento")
		}
	case "":
		tipoPersonaNull = true
		if len(newPersona.RazonSocial) > 0 {
			return clientErrorResponse("razonSocial debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.Nombre) > 0 {
			return clientErrorResponse("nombre debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.EstadoCUIT) > 0 {
			return clientErrorResponse("estadoCuit debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.Apellido) > 0 {
			return clientErrorResponse("apellido debe ser nulo cuando tipoPersona es nulo")
		}
		if newPersona.IDFormaJuridica != 0 {
			return clientErrorResponse("idFormaJuridica debe ser nulo cuando tipoPersona es nulo")
		}
		if newPersona.TipoDoc != 0 {
			return clientErrorResponse("tipoDoc debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.Documento) > 0 {
			return clientErrorResponse("documento debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.Sexo) > 0 {
			return clientErrorResponse("sexo debe ser nulo cuando tipoPersona es nulo")
		}
		if newPersona.MesCierre != 0 {
			return clientErrorResponse("mesCierre debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.FechaNacimiento) > 0 {
			return clientErrorResponse("fechaNacimiento debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.FechaFallecimiento) > 0 {
			return clientErrorResponse("fechaFallecimiento debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.FechaInscripcion) > 0 {
			return clientErrorResponse("fechaInscripcion debe ser nulo cuando tipoPersona es nulo")
		}
		if len(newPersona.FechaCierre) > 0 {
			return clientErrorResponse("fechaCierre debe ser nulo cuando tipoPersona es nulo")
		}
		if newPersona.NuevaCUIT != 0 {
			return clientErrorResponse("nuevaCuit debe ser nulo cuando tipoPersona es nulo")
		}
		if newPersona.Impuestos == nil || len(newPersona.Impuestos) == 0 {
			return clientErrorResponse("cuando tipoPerona es nulo debe informarse por lo menos un item en alguno de los arrarys (impuestos, actividades, domicilios, ...)")
		}
	default:
		return clientErrorResponse("tipoPersona [" + newPersona.TipoPersona + "] invalido, debe ser F (Humana), J (Juridica) o nulo (para indicar que no se aplican cambios sobre el asset topo Persona) ")
	}
	if !tipoPersonaNull {
		switch newPersona.EstadoCUIT {
		case "A":
		case "I":
		default:
			return clientErrorResponse("estadoCuit [" + newPersona.EstadoCUIT + "] invalido, debe ser A (Activa) o I (Inactiva)")
		}
	}

	key := getPersonaKey(newPersona)

	if exist, err := keyExists(APIstub, cuitStr); err == (Response{}) {
		return err
	} else if !exist {
		if tipoPersonaNull {
			return clientErrorResponse("No existe un asset [" + cuitStr + "] - Debe informarse los datos identificarios de la Persona")
		}
		log.Println("Putting [" + cuitStr + "]...")
		if err := APIstub.PutState(cuitStr, []byte("{}")); err != nil {
			return systemErrorResponse("Error putting cuitStr [" + cuitStr + "]: " + err.Error())
		}
	}

	var impuestos = newPersona.Impuestos
	newPersona.Impuestos = nil

	// Si tipo Persona es Null, no se guardan los datos de esa persona
	if !tipoPersonaNull {
		personaAsBytes, _ := json.Marshal(newPersona)

		log.Println("Putting [" + key + "]...")
		if err := APIstub.PutState(key, personaAsBytes); err != nil {
			return systemErrorResponse("Error putting key [" + key + "]: " + err.Error())
		}
	}

	if rows, err := s.commitPersonaImpuestos(APIstub, cuitStr, impuestos); err == (Response{}) {
		log.Println(err.Msg)
		return err
	} else {
		if !tipoPersonaNull {
			rows++
		}
		log.Println("Ok")
		return successResponse("Ok", rows)
	}
}
