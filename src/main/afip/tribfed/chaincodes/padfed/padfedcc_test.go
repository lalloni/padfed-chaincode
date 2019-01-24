package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getPesonaJSON(cuit uint64) string {
	var tipoPersona = "F"
	var razonSocial = ""
	var nombreApellido = `"nombre": "Pepe", "apellido": "Sanchez",`
	var idFormaJuridica = "0"
	var documento = `"documento": "27123456",`
	var tipoDoc = "1"
	var fechaNacimiento = `"fechaNacimiento":"1928-11-17",`
	if cuit >= 30000000000 {
		tipoPersona = "J"
		razonSocial = `"razonSocial":"THE BIRTH OF MARIA CONCETTA",`
		idFormaJuridica = "1"
		nombreApellido = ""
		tipoDoc = "0"
		documento = ""
		fechaNacimiento = ""
	}
	var personaJSON = `{
	"cuit":$cuit,"tipoPersona":"$tipoPersona","estadoCuit":"A",$razonSocial $nombreApellido
	"idFormaJuridica":$idFormaJuridica, "tipoDoc": $tipoDoc, $documento "fechaInscripcion":"1992-10-20","mesCierre":12, $fechaNacimiento
	"impuestos":[
		{"idImpuesto":30,"estado":"AC","periodo":199912},
		{"idImpuesto":217,"estado":"AC","periodo":199605},
		{"idImpuesto":301,"estado":"AC","periodo":199407},
		{"idImpuesto":103,"estado":"AC","periodo":201112}
	],
	"actividades":[
		{"codNomenclador":883,"idActividad":941100,"orden":1,"estado":"AC","periodo":201311}
	],
	"domicilios":[
		{"idTipoDomicilio":1,"orden":1,"idEstadoDomicilio":11,"idNomenclador":"401084021","codPostal":"1084","idProvincia":"0","localidad":"","calle":"DE MAYO AV.","numero":"568"},
		{"idTipoDomicilio":2,"orden":1,"idEstadoDomicilio":11,"idNomenclador":"401084021","codPostal":"1084","idProvincia":"0","localidad":"","calle":"DE MAYO AV.","numero":"568"}
	],
	"telefonos":[
		{"numero":"46788554","idEstadoTelefono":"3"}
	]
}`

	cuitStr := strconv.FormatUint(cuit, 10)
	personaJSON = strings.Replace(personaJSON, "$cuit", cuitStr, -1)
	personaJSON = strings.Replace(personaJSON, "$tipoPersona", tipoPersona, -1)
	personaJSON = strings.Replace(personaJSON, "$razonSocial", razonSocial, -1)
	personaJSON = strings.Replace(personaJSON, "$nombreApellido", nombreApellido, -1)
	personaJSON = strings.Replace(personaJSON, "$idFormaJuridica", idFormaJuridica, -1)
	personaJSON = strings.Replace(personaJSON, "$tipoDoc", tipoDoc, -1)
	personaJSON = strings.Replace(personaJSON, "$documento", documento, -1)
	personaJSON = strings.Replace(personaJSON, "$fechaNacimiento", fechaNacimiento, -1)
	return personaJSON
}

func TestValidPersonaJSON(t *testing.T) {
	var persona Persona
	var personaJSON = getPesonaJSON(30679638943)
	if err := argToPersona([]byte(personaJSON), &persona, JSON); err != (Response{}) {
		t.Error(err.Msg)
	}
	if getPersonaKey(&persona) != "PER_30679638943" {
		t.Error("Persona.Key no valida " + getPersonaKey(&persona))
	}
	if err := argToPersona([]byte("{error-dummy"), &persona, JSON); err != (Response{}) {
		t.Error("JSON invalido, debe dar error" + err.Msg)
	}
}

func TestCuit(t *testing.T) {
	_, err := getCUITArgs([]string{"1"})
	if err == nil {
		t.Error("Debe pinchar, el valor 1 no es un cuit valido.")
	}
	_, err = getCUITArgs([]string{"20066675573"})
	if err != nil {
		t.Error("No debe pinchar, el valor cuit es valido.")
	}
}

func TestValidImpuestosJSON(t *testing.T) {
	const cuit = 30679638943
	var impuestos Impuestos
	var personaJSON = getPesonaJSON(cuit)
	err := json.Unmarshal([]byte(personaJSON), &impuestos)

	if err != nil {
		t.Error("Error Failed to decode JSON of Impuestos")
	}

	if len(impuestos.Impuestos) != 4 {
		t.Error("Persona debe tener 4 impuestos y tiene " + strconv.Itoa(len(impuestos.Impuestos)))
	}
	if getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[0].IDImpuesto) != "PER_30679638943_IMP_30" {
		t.Error("1-Impuesto.Key no valido " + getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[0].IDImpuesto))
	}
	if getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[3].IDImpuesto) != "PER_30679638943_IMP_103" {
		t.Error("3-Impuesto.Key no valido " + getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[3].IDImpuesto))
	}
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func putPersona(t *testing.T, stub *shim.MockStub, cuit uint64) pb.Response {
	var personaJSON = getPesonaJSON(cuit)
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("putPersona"), []byte(cuitStr), []byte(personaJSON)})
}

func putPersonaProto(t *testing.T, stub *shim.MockStub, cuit uint64) pb.Response {
	var personaJSON = getPesonaJSON(cuit)
	var persona Persona

	argToPersona([]byte(personaJSON), &persona, JSON)

	personaPROTO, err := proto.Marshal(&persona)
	if err != nil {
		log.Fatal("marshaling to PROTOBUF error: ", err)
	}
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("putPersonaProto"), []byte(cuitStr), []byte(personaPROTO)})
}

func queryPersona(t *testing.T, stub *shim.MockStub, cuit uint64) pb.Response {
	cuitStr := strconv.FormatUint(cuit, 10)
	return stub.MockInvoke("1", [][]byte{[]byte("queryPersona"), []byte(cuitStr)})
}

func setInitTests(t *testing.T) *shim.MockStub {
	scc := new(SmartContract)
	stub := shim.NewMockStub("padfed", scc)
	scc.isModeTest = true
	checkInit(t, stub, [][]byte{})
	return stub
}

func TestPutPersona(t *testing.T) {
	stub := setInitTests(t)

	// Valid
	res := putPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}

	// Invalid cuit
	res = putPersona(t, stub, 1)
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con un cuit invalido debe dar error")
		t.FailNow()
	}

	// distinct cuits
	var personaJSON = getPesonaJSON(20255438795)
	res = stub.MockInvoke("1", [][]byte{[]byte("putPersona"), []byte("30679638943"), []byte(personaJSON)})
	if res.Status != shim.ERROR {
		fmt.Println("putPersona con cuits distintos debe dar error")
		t.FailNow()
	}
}

func TestPutPersonaProto(t *testing.T) {
	stub := setInitTests(t)
	// Valid
	res := putPersonaProto(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersonaProto", "cuit", "failed", string(res.Message))
		t.FailNow()
	} else {
		fmt.Println("putPersonaProto Ok!!!!")
	}
}

func TestPutPersonas(t *testing.T) {
	stub := setInitTests(t)

	var pJSON = `{"personas":[{"cuit":20066675573,"apellido":"GES","nombre":"THOMAS MICHAEL","tipoPersona":"F","estadoCuit":"A","tipoDoc":1,"documento":"6667557","sexo":"M","fechaNacimiento":"1928-11-17","impuestos":[{"idImpuesto":11,"estado":"BD","periodo":199901},{"idImpuesto":20,"estado":"BD","periodo":200907},{"idImpuesto":21,"estado":"BD","periodo":200907},{"idImpuesto":180,"estado":"AC","periodo":199807}],"categorias":[{"idCategoria":"11","estado":"BD","idImpuesto":20,"periodo":200907},{"idCategoria":"11","estado":"BD","idImpuesto":21,"periodo":200907}],"actividades":[{"codNomenclador":883,"idActividad":692000,"orden":1,"estado":"AC","periodo":201311}],"domicilios":[{"idTipoDomicilio":1,"orden":1,"idEstadoDomicilio":2,"idNomenclador":"3541","codPostal":"5891","idProvincia":"3","localidad":"VILLA CURA BROCHERO","calle":"HIPOLITO IRIGOYEN","numero":"57"},{"idTipoDomicilio":2,"orden":1,"idEstadoDomicilio":9,"idNomenclador":"3541","codPostal":"5891","idProvincia":"3","localidad":"VILLA CURA BROCHERO","calle":"SAN MARTIN ESQ IRIGO","numero":"8"}]},{"cuit":20066758193,"apellido":"RACCONTARE","nombre":"GUSTAVO FABIAN","tipoPersona":"F","estadoCuit":"A","tipoDoc":1,"documento":"6675819","sexo":"M","fechaNacimiento":"1933-01-22","impuestos":[{"idImpuesto":11,"estado":"AC","periodo":190101},{"idImpuesto":30,"estado":"AC","periodo":200408},{"idImpuesto":32,"estado":"BD","periodo":200408},{"idImpuesto":180,"estado":"AC","periodo":199105},{"idImpuesto":301,"estado":"AC","periodo":199407},{"idImpuesto":308,"estado":"AC","periodo":196501}],"categorias":[{"idCategoria":"501","estado":"AC","idImpuesto":308,"periodo":200703}],"actividadList":[],"domicilios":[{"idTipoDomicilio":1,"orden":1,"idEstadoDomicilio":6,"idNomenclador":"6024","codPostal":"3315","idProvincia":"19","localidad":"LEANDRO N. ALEM","calle":"RIVADAVIA","numero":"572"},{"idTipoDomicilio":2,"orden":1,"idEstadoDomicilio":1,"idNomenclador":"6024","codPostal":"3315","idProvincia":"19","localidad":"LEANDRO N. ALEM","calle":"URUGUAY","numero":"287"}]}]}`
	res := stub.MockInvoke("1", [][]byte{[]byte("putPersonas"), []byte(pJSON)})
	if res.Status != shim.OK {
		fmt.Println("putPersonas", string(res.Message))
		t.FailNow()
	} else {
		fmt.Println("putPersonas Ok!!!!")
	}

}

func TestPutPersonasProto(t *testing.T) {
	stub := setInitTests(t)

	cuils := []uint64{20066042333, 20066675573, 20066806163, 20068854785, 20176058650}
	personas := Personas{}
	var p []*Persona

	for _, cuit := range cuils {
		var persona Persona
		json.Unmarshal([]byte(getPesonaJSON(cuit)), &persona)
		p = append(p, &persona)
	}
	personas.Personas = p

	personasPROTO, err := proto.Marshal(&personas)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	res := stub.MockInvoke("1", [][]byte{[]byte("putPersonasProto"), []byte(personasPROTO)})
	if res.Status != shim.OK {
		fmt.Println("putPersonaProto", "cuit", "failed", string(res.Message))
		t.FailNow()
	} else {
		fmt.Println("putPersonasProto Ok!!!!")
	}
}
func TestQueryPersona(t *testing.T) {
	stub := setInitTests(t)
	res := putPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
	res = queryPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("queryPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
	fmt.Println("queryPersona ", string(res.Payload))

}

func TestPutPersonaImpuestos(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}

	impuestosJSON := `{"impuestos":[{"idImpuesto":30,"estado":"AC","periodo":199912},{"idImpuesto":31,"idOrg":901,"estado":"AC","periodo":199912}]}`

	res = stub.MockInvoke("1", [][]byte{[]byte("putPersonaImpuestos"), []byte("30679638943"), []byte(impuestosJSON)})
	if res.Status != shim.OK {
		fmt.Println("putPersonaImpuestos error", string(res.Message))
		t.FailNow()
	}
}

func TestCreateTxConfirmable(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 20255438795)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("1", [][]byte{
		[]byte("createTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("2"),
		[]byte("1"),
		[]byte("Impuesto"),
		[]byte(`{"idImpuesto":217,"estado":"B","periodo":199605}`)})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("1", [][]byte{
		[]byte("createTxConfirmable"),
		[]byte("20255438795"),
		[]byte("2"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("2"),
		[]byte("1"),
		[]byte("Impuesto"),
		[]byte(`{"idImpuesto":217,"estado":"B","periodo":199605}`)})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}
}

func TestResponseTxConfirmable(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 20255438795)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
	// createTxConfirmable
	res = stub.MockInvoke("1", [][]byte{
		[]byte("createTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("2"),
		[]byte("1"),
		[]byte("Impuesto"),
		[]byte(`{"idImpuesto":217,"estado":"B","periodo":199605}`)})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}
	// responseTxConfirmable -> confirma cambio
	res = stub.MockInvoke("1", [][]byte{
		[]byte("responseTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("1")})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}
}

/*func TestQueryHistory(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 20255438795)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
	// createTxConfirmable
	res = stub.MockInvoke("1", [][]byte{
		[]byte("createTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("2"),
		[]byte("1"),
		[]byte("Impuesto"),
		[]byte(`{"idImpuesto":217,"estado":"B","periodo":199605}`)})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}
	// responseTxConfirmable -> confirma cambio
	res = stub.MockInvoke("1", [][]byte{
		[]byte("responseTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("1")})
	if res.Status != shim.OK {
		fmt.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}

	// responseTxConfirmable -> confirma cambio
	res = stub.MockInvoke("1", [][]byte{
		[]byte("queryHistory"),
		[]byte("PER_20255438795_IMP_217"),
	})
	if res.Status != shim.OK {
		fmt.Println("queryHistory error", string(res.Message))
		t.FailNow()
	}
}*/

func TestQueryTxConfirmables(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 20255438795)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
	// createTxConfirmable
	res = stub.MockInvoke("1", [][]byte{
		[]byte("createTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("2"),
		[]byte("900"),
		[]byte("Impuesto"),
		[]byte(`{"idImpuesto":217,"estado":"B","periodo":199605}`)})
	if res.Status != shim.OK {
		log.Println("createTxConfirmable error", string(res.Message))
		t.FailNow()
	}
	// responseTxConfirmable -> confirma cambio
	res = stub.MockInvoke("1", [][]byte{
		[]byte("responseTxConfirmable"),
		[]byte("20255438795"),
		[]byte("1"),
		[]byte("2"),
		[]byte("900"),
		[]byte("2002-10-02T15:00:00.05Z"),
		[]byte("1")})
	// query ->
	res = stub.MockInvoke("1", [][]byte{
		[]byte("queryTxConfirmables"),
		[]byte("900"),
		[]byte("1"),
		[]byte("")})

	if res.Status != shim.OK {
		log.Println("queryTxConfirmables error" + string(res.Message))
		t.FailNow()
	}
}

func TestDelPersonasByRange(t *testing.T) {
	stub := setInitTests(t)

	cuils := []uint64{20066042333, 20066675573, 20066806163, 20068854785, 20176058650}
	// creo varias personas
	for _, cuil := range cuils {
		res := putPersona(t, stub, cuil)
		if res.Status != shim.OK {
			fmt.Println("putPersona", "cuit", "failed", string(res.Message))
			t.FailNow()
		}
	}

	res := stub.MockInvoke("1", [][]byte{[]byte("delPersonasByRange"), []byte("20066600000"), []byte("20068900000")})
	if res.Status != shim.OK {
		fmt.Println("delPersonasByRange", "cuit", "failed", string(res.Message))
		t.FailNow()
	}

	fmt.Println("--- ESTADO FINAL ---")
	res = stub.MockInvoke("1", [][]byte{[]byte("queryAllPersona")})
	if res.Status != shim.OK {
		fmt.Println("queryAllPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}
}
