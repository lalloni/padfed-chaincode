package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var impuestosPrecarga = `[                                                                              
    {"impuesto":20,"idOrg":1,"nombre":"MONOTRIBUTO"},                                                                                                 
    {"impuesto":21,"idOrg":1,"nombre":"MONOTRIBUTO AUTONOMO"},                                                                                        
    {"impuesto":22,"idOrg":1,"nombre":"MONOTRIBUTO SEG.SOCIAL"},                                                                                      
    {"impuesto":23,"idOrg":1,"nombre":"MONOTRIBUTO-INTEG.DE SOCIEDAD"},                                                                           
    {"impuesto":30,"idOrg":1,"nombre":"IVA"},                                                                                        
    {"impuesto":32,"idOrg":1,"nombre":"IVA EXENTO"},                                                                                                  
    {"impuesto":33,"idOrg":1,"nombre":"IVA RESPONSABLE NO INSCRIPTO"},                                                                                
    {"impuesto":34,"idOrg":1,"nombre":"IVA NO ALCANZADO"},                                                                                    
    {"impuesto":301,"idOrg":1,"nombre":"EMPLEADOR-APORTES SEG. SOCIAL"},                                                                            
    {"impuesto":308,"idOrg":1,"nombre":"APORTES SEG.SOCIAL AUTONOMOS"},    
    {"impuesto":5243,"idOrg":904,"nombre":"INGRESOS BRUTOS CORDOBA"},                                                                                 
    {"impuesto":5244,"idOrg":904,"nombre":"CONTRIBUCION MUNICIPAL"}                                                                         
    ]`

func (s *SmartContract) setInitImpuestos(APIstub shim.ChaincodeStubInterface) Response {
	return s.putParamImpuestos(APIstub, []string{impuestosPrecarga})
}
