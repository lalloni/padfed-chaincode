package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getPesonaJSON(cuit uint64) string {
	var tipoPer = "F"
	var razonSocial = ""
	var nombreApellido = `"nombre": "Pepe", "apellido": "Sanchez",`
	var formaJuridica = "0"
	var doc = `"doc": "27123456",`
	var tipoDoc = "1"
	var nacimiento = `"nacimiento":"1928-11-17",`
	if cuit >= 30000000000 {
		tipoPer = "J"
		razonSocial = `"razonSocial":"THE BIRTH OF MARIA CONCETTA",`
		formaJuridica = "1"
		nombreApellido = ""
		tipoDoc = "0"
		doc = ""
		nacimiento = ""
	}
	var personaJSON = `{
	"cuit":$cuit,"tipo":"$tipoPer","estado":"A",$razonSocial $nombreApellido
	"formaJuridica":$formaJuridica, "tipoDoc":$tipoDoc, $doc "inscripcion":"1992-10-20","mesCierre":12, $nacimiento
	"impuestos":[
		{"impuesto":30,"estado":"AC","periodo":199912},
		{"impuesto":5100,"estado":"AC","periodo":199605},
		{"impuesto":301,"estado":"AC","periodo":199407},
		{"impuesto":34,"estado":"AC","periodo":201112}
	],
	"actividades":[
		{"nomenclador":883,"id":941100,"orden":1,"estado":"AC","periodo":201311}
	],
	"domicilios":[
		{"tipo":1,"orden":1,"estado":11,"nomenclador":401084021,"codPostal":"1084","provincia":0,"localidad":"","calle":"DE MAYO AV.","numero":"568"},
		{"tipo":2,"orden":1,"estado":11,"nomenclador":401084021,"codPostal":"1084","provincia":0,"localidad":"","calle":"DE MAYO AV.","numero":"568"}
	],
	"telefonos":[
		{"numero":"46788554","idEstadoTelefono":"3"}
	]
}`

	cuitStr := strconv.FormatUint(cuit, 10)
	personaJSON = strings.Replace(personaJSON, "$cuit", cuitStr, -1)
	personaJSON = strings.Replace(personaJSON, "$tipoPer", tipoPer, -1)
	personaJSON = strings.Replace(personaJSON, "$razonSocial", razonSocial, -1)
	personaJSON = strings.Replace(personaJSON, "$nombreApellido", nombreApellido, -1)
	personaJSON = strings.Replace(personaJSON, "$formaJuridica", formaJuridica, -1)
	personaJSON = strings.Replace(personaJSON, "$tipoDoc", tipoDoc, -1)
	personaJSON = strings.Replace(personaJSON, "$doc", doc, -1)
	personaJSON = strings.Replace(personaJSON, "$nacimiento", nacimiento, -1)
	return personaJSON
}

func TestValidPersona(t *testing.T) {
	var persona Persona
	var personaJSON = getPesonaJSON(30679638943)
	if err := argToPersona([]byte(personaJSON), &persona); err.isError() {
		t.Error(err.Msg)
	}
	if getPersonaKey(&persona) != "PER_30679638943" {
		t.Error("Persona.Key no valida " + getPersonaKey(&persona))
	}
	if err := argToPersona([]byte("{error-dummy"), &persona); err.isOk() {
		t.Error("JSON invalido, debe dar error " + err.Msg)
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

func TestValimpuestosJSON(t *testing.T) {
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
	if getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[0].Impuesto) != "PER_30679638943_IMP_30" {
		t.Error("1-Impuesto.Key no valido " + getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[0].Impuesto))
	}
	if getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[3].Impuesto) != "PER_30679638943_IMP_34" {
		t.Error("3-Impuesto.Key no valido " + getImpuestoKeyByCuitId(cuit, impuestos.Impuestos[3].Impuesto))
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

func TestPutPersonas(t *testing.T) {
	stub := setInitTests(t)

	var pJSON = `{"personas":[{"cuit":20066675573,"apellido":"GES","nombre":"THOMAS MICHAEL","tipo":"F","estado":"A","tipoDoc":1,"doc":"6667557","sexo":"M","nacimiento":"1928-11-17","impuestos":[{"impuesto":20,"estado":"BD","periodo":199901},{"impuesto":5243,"estado":"BD","periodo":200907},{"impuesto":21,"estado":"BD","periodo":200907},{"impuesto":5244,"estado":"AC","periodo":199807}],"categorias":[{"idCategoria":"11","estado":"BD","impuesto":20,"periodo":200907},{"idCategoria":"11","estado":"BD","impuesto":23,"periodo":200907}],"actividades":[{"nomenclador":883,"id":692000,"orden":1,"estado":"AC","periodo":201311}],"domicilios":[{"tipo":1,"orden":1,"idEstadoDomicilio":2,"nomenclador":3541,"codPostal":"5891","provincia":3,"localidad":"VILLA CURA BROCHERO","calle":"HIPOLITO IRIGOYEN","numero":"57"},{"tipo":2,"orden":1,"idEstadoDomicilio":9,"nomenclador":3541,"codPostal":"5891","provincia":3,"localidad":"VILLA CURA BROCHERO","calle":"SAN MARTIN ESQ IRIGO","numero":"8"}]},{"cuit":20066758193,"apellido":"RACCONTARE","nombre":"GUSTAVO FABIAN","tipo":"F","estado":"A","tipoDoc":1,"doc":"6675819","sexo":"M","nacimiento":"1933-01-22","impuestos":[{"impuesto":20,"estado":"AC","periodo":190101},{"impuesto":30,"estado":"AC","periodo":200408},{"impuesto":32,"estado":"BD","periodo":200408},{"impuesto":5244,"estado":"AC","periodo":199105},{"impuesto":301,"estado":"AC","periodo":199407},{"impuesto":5100,"estado":"AC","periodo":196501}],"categorias":[{"idCategoria":"501","estado":"AC","impuesto":5100,"periodo":200703}],"actividadList":[],"domicilios":[{"tipo":1,"orden":1,"idEstadoDomicilio":6,"nomenclador":6024,"codPostal":"3315","provincia":19,"localidad":"LEANDRO N. ALEM","calle":"RIVADAVIA","numero":"572"},{"tipo":2,"orden":1,"idEstadoDomicilio":1,"nomenclador":6024,"codPostal":"3315","provincia":19,"localidad":"LEANDRO N. ALEM","calle":"URUGUAY","numero":"287"}]}]}`
	res := stub.MockInvoke("1", [][]byte{[]byte("putPersonas"), []byte(pJSON)})
	if res.Status != shim.OK {
		fmt.Println("putPersonas", string(res.Message))
		t.FailNow()
	} else {
		fmt.Println("putPersonas Ok!!!!")
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

	impuestosJSON := `{"impuestos":[{"impuesto":30,"estado":"AC","periodo":199912},{"impuesto":32,"estado":"AC","periodo":199912}]}`

	res = stub.MockInvoke("1", [][]byte{[]byte("putPersonaImpuestos"), []byte("30679638943"), []byte(impuestosJSON)})
	if res.Status != shim.OK {
		fmt.Println("putPersonaImpuestos error", string(res.Message))
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

func TestDelPersonaAssets(t *testing.T) {
	stub := setInitTests(t)

	res := putPersona(t, stub, 30679638943)
	if res.Status != shim.OK {
		fmt.Println("putPersona", "cuit", "failed", string(res.Message))
		t.FailNow()
	}

	impuestosToDel := `["PER_30679638943_IMP_30","PER_30679638943_IMP_124"]`

	res = stub.MockInvoke("1", [][]byte{[]byte("delPersonaAssets"), []byte("30679638943"), []byte(impuestosToDel)})
	if res.Status != shim.OK {
		fmt.Println("delPersonaAssets error", string(res.Message))
		t.FailNow()
	}
}
