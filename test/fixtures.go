package test

import (
	"strconv"
	"strings"
)

func GetPersonaJSON(cuit uint64) string {
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
