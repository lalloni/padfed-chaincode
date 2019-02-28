package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var impuestosPrecarga = `[
    {"impuesto":10,"idOrg":1,"nombre":"GANANCIAS SOCIEDADES"},                                                                                        
    {"impuesto":11,"idOrg":1,"nombre":"GANANCIAS PERSONAS FISICAS"},                                                                            
    {"impuesto":20,"idOrg":1,"nombre":"MONOTRIBUTO"},                                                                                                 
    {"impuesto":21,"idOrg":1,"nombre":"MONOTRIBUTO AUTONOMO"},                                                                                        
    {"impuesto":22,"idOrg":1,"nombre":"MONOTRIBUTO SEG.SOCIAL"},                                                                                      
    {"impuesto":23,"idOrg":1,"nombre":"MONOTRIBUTO-INTEG.DE SOCIEDAD"},                                                                           
    {"impuesto":30,"idOrg":1,"nombre":"IVA"},                                                                                        
    {"impuesto":32,"idOrg":1,"nombre":"IVA EXENTO"},                                                                                                  
    {"impuesto":33,"idOrg":1,"nombre":"IVA RESPONSABLE NO INSCRIPTO"},                                                                                
    {"impuesto":34,"idOrg":1,"nombre":"IVA NO ALCANZADO"},
    {"impuesto":80,"idOrg":902,"nombre":"IIBB - PCIA BUENOS AIRES"},
    {"impuesto":301,"idOrg":1,"nombre":"EMPLEADOR-APORTES SEG. SOCIAL"},
    {"impuesto":5100,"idOrg":901,"nombre":"IIBB - CABA"},                                                                                 
    {"impuesto":5200,"idOrg":904,"nombre":"IIBB - CORDOBA"},    
    {"impuesto":5243,"idOrg":904,"nombre":"INGRESOS BRUTOS CORDOBA"},                                                                                 
    {"impuesto":5244,"idOrg":904,"nombre":"CONTRIBUCION MUNICIPAL"},
    {"impuesto":5300,"idOrg":903,"nombre":"IIBB - CATAMARCA"},                                                                                 
    {"impuesto":5800,"idOrg":900,"nombre":"CONVENIO MULTILATERAL"}                                                                             
    ]`

func (s *SmartContract) setInitImpuestos(APIstub shim.ChaincodeStubInterface) Response {
	return s.putParamImpuestos(APIstub, []string{impuestosPrecarga})
}
