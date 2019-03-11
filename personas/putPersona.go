package personas

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/fabric"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/inscripciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/model"
)

func PutPersona(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 2 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se esperan 2 (CUIT, PERSONA)")
	}
	var cuit uint64
	var err error
	var res *fabric.Response
	cuitStr := args[0]
	log.Print("cuit recibido [" + cuitStr + "]")
	if cuit, err = helpers.GetCUIT(cuitStr); err != nil {
		return fabric.ClientErrorResponse("CUIT [" + cuitStr + "] invalido")
	}
	newPersona := &model.Persona{}
	if res = ArgToPersona([]byte(args[1]), newPersona); !res.IsOK() {
		return res
	}
	return SavePersona(stub, cuit, newPersona)
}

func PutPersonas(stub shim.ChaincodeStubInterface, args []string) *fabric.Response {
	if len(args) != 1 {
		return fabric.ClientErrorResponse("Número incorrecto de argumentos. Se espera 1 (PERSONAS)")
	}
	newPersonas := &model.Personas{}
	if err := ArgToPersonas([]byte(args[0]), newPersonas); !err.IsOK() {
		return err
	}

	rows := 0
	for _, p := range newPersonas.Personas {
		log.Printf("Grabando persona %d", p.CUIT)
		res := SavePersona(stub, p.CUIT, p)
		if res.Status != shim.OK {
			res.WrongItem = rows
			return res
		}
		rows++
	}
	log.Println(strconv.Itoa(rows) + " personas processed !!!")
	return fabric.SuccessResponse("Ok", rows)
}

func SavePersona(stub shim.ChaincodeStubInterface, cuit uint64, p *model.Persona) *fabric.Response {

	cuits := strconv.FormatUint(cuit, 10)
	if cuit != p.CUIT {
		return fabric.ClientErrorResponse("El parametro cuit [" + cuits + "] y la cuit [" + strconv.FormatUint(p.CUIT, 10) + "] en la Persona deben ser iguales")
	}

	switch p.Tipo {
	case "F":
		if res := checkPersonaFisica(cuits, p); res != nil {
			return res
		}
		if res := checkEstadoPersona(p); res != nil {
			return res
		}
	case "J":
		if res := checkPersonaJuridica(cuits, p); res != nil {
			return res
		}
		if res := checkEstadoPersona(p); res != nil {
			return res
		}
	case "":
		if res := checkPersonaSinTipo(p); res != nil {
			return res
		}
	default:
		return fabric.ClientErrorResponse("tipo [" + p.Tipo + "] invalido, debe ser F (Humana), J (Juridica) o nulo (para indicar que no se aplican cambios sobre el asset topo Persona) ")
	}

	key := GetPersonaKey(p)

	exist, err := fabric.KeyExists(stub, cuits)
	if !err.IsOK() {
		return err
	}

	if !exist {
		if p.Tipo == "" {
			return fabric.ClientErrorResponse("No existe un asset [" + cuits + "] - Debe informarse los datos identificarios de la Persona")
		}
		log.Println("Putting [" + cuits + "]...")
		if err := stub.PutState(cuits, []byte("{}")); err != nil {
			return fabric.SystemErrorResponse("Error putting cuitStr [" + cuits + "]: " + err.Error())
		}
	}

	var impuestos = p.Impuestos
	p.Impuestos = nil

	// Sólo se guardan los datos si la persona tiene tipo
	if p.Tipo != "" {
		personaAsBytes, _ := json.Marshal(p)

		log.Println("Putting [" + key + "]...")
		if err := stub.PutState(key, personaAsBytes); err != nil {
			return fabric.SystemErrorResponse("Error putting key [" + key + "]: " + err.Error())
		}
	}

	rows, err := inscripciones.CommitPersonaImpuestos(stub, cuits, impuestos)
	if !err.IsOK() {
		return err
	}

	if p.Tipo != "" {
		rows++
	}

	log.Print("Ok")

	return fabric.SuccessResponse("Ok", rows)
}

func checkEstadoPersona(p *model.Persona) *fabric.Response {
	if p.Estado != "A" && p.Estado != "I" {
		return fabric.ClientErrorResponse("estado [" + p.Estado + "] invalido, debe ser A (Activa) o I (Inactiva)")
	}
	return nil
}

func checkPersonaFisica(cuit string, p *model.Persona) *fabric.Response {
	if !strings.HasPrefix(cuit, "2") {
		return fabric.ClientErrorResponse("la cuit [" + cuit + "] debe comenzar con 2 cuando corresponde a una persona humana")
	}
	p.Nombre = strings.ToUpper(strings.Trim(p.Nombre, " "))
	p.Apellido = strings.ToUpper(strings.Trim(p.Apellido, " "))
	if len(p.Apellido) == 0 {
		return fabric.ClientErrorResponse("persona humana sin apellido")
	}
	if len(p.RazonSocial) > 0 {
		return fabric.ClientErrorResponse("persona humana con razonSocial [" + p.RazonSocial + "]")
	}
	if p.FormaJuridica != 0 {
		return fabric.ClientErrorResponse("persona humana con formaJuridica")
	}
	if !(p.TipoDoc >= 1 && p.TipoDoc <= 99) {
		return fabric.ClientErrorResponse("tipoDoc [" + strconv.Itoa(int(p.TipoDoc)) + "] invalido, debe ser un entero entre 1 y 99")
	}
	if !(p.MesCierre == 0 || p.MesCierre == 12) {
		return fabric.ClientErrorResponse("mesCierre [" + strconv.Itoa(int(p.MesCierre)) + "] invalido. Para una persona fisica debe ser 12.")
	}
	p.Doc = strings.ToUpper(strings.Trim(p.Doc, " "))
	if len(p.Doc) == 0 {
		return fabric.ClientErrorResponse("persona humana sin doc")
	}
	if len(p.Nacimiento) == 0 {
		return fabric.ClientErrorResponse("persona humana sin nacimiento")
	}
	return nil
}

func checkPersonaJuridica(cuit string, p *model.Persona) *fabric.Response {
	if !strings.HasPrefix(cuit, "3") {
		return fabric.ClientErrorResponse("la cuit [" + cuit + "] debe comenzar con 3 cuando corresponde a una persona juridica")
	}
	p.RazonSocial = strings.ToUpper(strings.Trim(p.RazonSocial, " "))
	if len(p.RazonSocial) == 0 {
		return fabric.ClientErrorResponse("persona juridica sin razonSocial")
	}
	if len(p.Nombre) > 0 {
		return fabric.ClientErrorResponse("persona juridica con nombre [" + p.Nombre + "], solo debe tener RazonSocial")
	}
	if len(p.Apellido) > 0 {
		return fabric.ClientErrorResponse("persona juridica con apellido [" + p.Apellido + "], solo debe tener RazonSocial")
	}
	if len(p.Materno) > 0 {
		return fabric.ClientErrorResponse("persona juridica con materno [" + p.Materno + "], solo debe tener RazonSocial")
	}
	if len(p.Sexo) > 0 {
		return fabric.ClientErrorResponse("persona juridica con sexo [" + p.Sexo + "], solo debe tener RazonSocial")
	}
	if !(p.FormaJuridica >= 0 && p.FormaJuridica <= 999) {
		return fabric.ClientErrorResponse("formaJuridica [" + strconv.Itoa(int(p.FormaJuridica)) + "] invalida")
	}
	if !(p.MesCierre >= 1 && p.MesCierre <= 12) {
		return fabric.ClientErrorResponse("mesCierre [" + strconv.Itoa(int(p.MesCierre)) + "] invalido. Para una persona juridica debe ser un entero entre 1 y 12.")
	}
	if p.TipoDoc != 0 {
		return fabric.ClientErrorResponse("persona juridica con tipoDoc")
	}
	if len(p.Doc) > 0 {
		return fabric.ClientErrorResponse("persona juridica con doc")
	}
	if len(p.Nacimiento) > 0 {
		return fabric.ClientErrorResponse("persona juridica con nacimiento")
	}
	if len(p.Fallecimiento) > 0 {
		return fabric.ClientErrorResponse("persona juridica con fallecimiento")
	}
	return nil
}

func checkPersonaSinTipo(p *model.Persona) *fabric.Response {
	if len(p.RazonSocial) > 0 {
		return fabric.ClientErrorResponse("razonSocial debe ser nulo cuando tipo es nulo")
	}
	if len(p.Nombre) > 0 {
		return fabric.ClientErrorResponse("nombre debe ser nulo cuando tipo es nulo")
	}
	if len(p.Estado) > 0 {
		return fabric.ClientErrorResponse("estado debe ser nulo cuando tipo es nulo")
	}
	if len(p.Apellido) > 0 {
		return fabric.ClientErrorResponse("apellido debe ser nulo cuando tipo es nulo")
	}
	if len(p.Materno) > 0 {
		return fabric.ClientErrorResponse("materno debe ser nulo cuando tipo es nulo")
	}
	if p.FormaJuridica != 0 {
		return fabric.ClientErrorResponse("formaJuridica debe ser nulo cuando tipo es nulo")
	}
	if p.TipoDoc != 0 {
		return fabric.ClientErrorResponse("tipoDoc debe ser nulo cuando tipo es nulo")
	}
	if len(p.Doc) > 0 {
		return fabric.ClientErrorResponse("doc debe ser nulo cuando tipo es nulo")
	}
	if len(p.Sexo) > 0 {
		return fabric.ClientErrorResponse("sexo debe ser nulo cuando tipo es nulo")
	}
	if len(p.Pais) > 0 {
		return fabric.ClientErrorResponse("pais debe ser nulo cuando tipo es nulo")
	}
	if p.MesCierre != 0 {
		return fabric.ClientErrorResponse("mesCierre debe ser nulo cuando tipo es nulo")
	}
	if len(p.Nacimiento) > 0 {
		return fabric.ClientErrorResponse("nacimiento debe ser nulo cuando tipo es nulo")
	}
	if len(p.Fallecimiento) > 0 {
		return fabric.ClientErrorResponse("fallecimiento debe ser nulo cuando tipo es nulo")
	}
	if len(p.Inscripcion) > 0 {
		return fabric.ClientErrorResponse("inscripcion debe ser nulo cuando tipo es nulo")
	}
	if len(p.FechaCierre) > 0 {
		return fabric.ClientErrorResponse("fechaCierre debe ser nulo cuando tipo es nulo")
	}
	if p.NuevaCUIT != 0 {
		return fabric.ClientErrorResponse("nuevaCuit debe ser nulo cuando tipo es nulo")
	}
	if p.Impuestos == nil || len(p.Impuestos) == 0 {
		return fabric.ClientErrorResponse("cuando tipo es nulo debe informarse por lo menos un item en alguno de los arrarys (impuestos, actividades, domicilios, ...)")
	}
	return nil
}
