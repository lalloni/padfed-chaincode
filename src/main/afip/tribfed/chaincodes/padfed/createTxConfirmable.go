package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
createTxConfirmable Crear una transacción confirmable.
Parámetros:
	1- p_cuit CUIT
	2- p_idTxc Entero entre 1 y 10000M
	3- p_fechahoraTxc yyyy-MM-dd'T'HH:mm:ss.SSSX
	4- p_tipoTxc String con valores:
		1: “CM – ALTA EN JURISDICCION”
		2: “CM – CAMBIO DE ESTADO EN JURISDICCION”
	5- p_idOrg Entero entre 900 y 999
	6- p_assetType “Impuesto”
	7- p_assetValue String json

*/
func (s *SmartContract) createTxConfirmable(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 7 {
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba 7  parámetros  con  {P_CUIT, P_IDTXC, P_FECHAHORATXC, P_TIPOTXC, P_IDORG, P_ASSETTYPE, P_ASSETVALUE}")
	}

	cuit, err := getCUITArgs(args)
	if err != nil {
		return clientErrorResponse(err.Error())
	}
	var idTxc int64
	if idTxc, err = strconv.ParseInt(args[1], 10, 64); err != nil {
		return clientErrorResponse("idTxc debe ser un número " + args[1])
	}

	var fechahoraTxc time.Time
	if fechahoraTxc, err = time.Parse(time.RFC3339, args[2]); err != nil {
		return clientErrorResponse("fechahoraTxc debe ser una fecha con formato yyyy-MM-dd'T'HH:mm:ss.SSSX " + args[2])
	}
	var tipoTxc int
	if tipoTxc, err = strconv.Atoi(args[3]); err != nil || (tipoTxc != 1 && tipoTxc != 2) {
		return clientErrorResponse("tipoTxc debe ser un número con valor (1 ó 2)")
	}

	var idOrg int
	if idOrg, err = strconv.Atoi(args[4]); err != nil {
		return clientErrorResponse("idOrg debe ser un número " + args[4])
	}

	assetType := args[5]
	if assetType != "Impuesto" {
		return clientErrorResponse("assetType debe ser Impuesto " + args[5])
	}

	var assetValue Impuesto
	if err = json.Unmarshal([]byte(args[6]), &assetValue); err != nil {
		return systemErrorResponse(err.Error())
	}

	// Persona
	exists, errP := keyExists(APIstub, strconv.FormatUint(cuit, 10))
	if errP.isError() {
		return errP
	} else if !exists {
		return clientErrorResponse("CUIT [" + strconv.FormatUint(cuit, 10) + "] inexistente")
	}

	// Validaciones
	txc, _, errT := findTXConfirmable(APIstub, idOrg, uint64(idTxc))
	if errT == nil {
		return clientErrorResponse("Ya existe una TXConfirmable con id " + args[1])
	}

	if assetValue.Estado != "A" && assetValue.Estado != "B" && assetValue.Estado != "E" {
		return clientErrorResponse("estado debe ser A,B ó E")
	}

	if tipoTxc == 1 && strings.HasPrefix(assetValue.Estado, "B") {
		return clientErrorResponse("No se puede crear una TXC de tipo " + strconv.Itoa(tipoTxc) + " asignadole estado " + assetValue.Estado)
	}

	var impuesto Impuesto
	var existImpuesto bool
	if impuesto, _, err = findImpuesto(APIstub, cuit, assetValue.IDImpuesto); err == nil {
		existImpuesto = true
	}

	if tipoTxc == 2 && !existImpuesto {
		return clientErrorResponse("No se puede crear una TXC de tipo " + strconv.Itoa(tipoTxc) + " con key PER_" + strconv.FormatUint(cuit, 10) + "_IMP_" + strconv.Itoa(int(assetValue.IDImpuesto)) + " porque no existe un asset con esa key")
	}

	if tipoTxc == 2 && existImpuesto && strings.HasPrefix(impuesto.Estado, assetValue.Estado) {
		return clientErrorResponse("No se puede crear una TXC de tipo " + strconv.Itoa(tipoTxc) + " con key PER_" + strconv.FormatUint(cuit, 10) + "_IMP_" + strconv.Itoa(int(assetValue.IDImpuesto)) + " porque existe un asset con esa misma key y con el mismo estado" + assetValue.Estado)
	}

	txc.CUIT = cuit
	txc.IDTxc = uint64(idTxc)

	if existImpuesto {
		impuesto.IDTxc = txc.IDTxc
		key := getImpuestoKeyByCuitId(cuit, impuesto.IDImpuesto)
		impuestoAsBytes, _ := json.Marshal(impuesto)
		log.Print("Se actualiza asset con key " + key)
		log.Print("AssetValue " + string(impuestoAsBytes))
		if err = APIstub.PutState(key, impuestoAsBytes); err != nil {
			return systemErrorResponse("Error al guardar Impuesto - " + key + ", error: " + err.Error())
		}
	}

	txc.FechaHoraTxc = fechahoraTxc
	txc.TipoTxc = tipoTxc
	txc.IDOrganismo = idOrg
	txc.AssetType = assetType
	assetValueAsBytes, _ := json.Marshal(assetValue)
	txc.AssetValue = string(assetValueAsBytes)

	txcAsBytes, _ := json.Marshal(txc)
	fmt.Println("Creando TXConfirmable " + getTxConfirmableKey(txc.IDOrganismo, txc.IDTxc) + ", assetValue: " + txc.AssetValue)
	if err = APIstub.PutState(getTxConfirmableKey(txc.IDOrganismo, txc.IDTxc), txcAsBytes); err != nil {
		return systemErrorResponse("Error al guardar TXConfirmable - " + err.Error())
	}
	return successResponse("", 0)
}
