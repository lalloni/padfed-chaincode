package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	TIPO_RESPUESTA_CONFIRMADO = 1
	TIPO_RESPUESTA_RECHAZADO  = 2
)

/*
responseTxConfirmable Permitir que un Organismo confirme o rechace un transacción confirmable.
Parámetros:
	1- p_cuit	CUIT
	2- p_idTxc	Entero entre 1 y 10000M
	3- p_tipoTxc	Entero entre 1 y 2
	4- p_idOrg 	Entero entre 900 y 999
	5- p_fechahoraRespuesta 	yyyy-MM-dd'T'HH:mm:ss.SSSX
	6- p_tipoRespuesta 	1:Confirmado | 2:Rechazado
*/
func (s *SmartContract) responseTxConfirmable(APIstub shim.ChaincodeStubInterface, args []string) Response {
	if len(args) != 6 {
		return clientErrorResponse("Número incorrecto de parámetros. Se esperaba 6  parámetros  con  {p_cuit, p_idTxc, p_tipoTxc, p_idOrg, p_fechahoraRespuesta, p_tipoRespuesta}")
	}
	var err error
	//cuit, err := getCUITArgs(args)
	//if err != nil {
	//	return s.businessErrorResponse(err)
	//}
	var idTxc int64
	if idTxc, err = strconv.ParseInt(args[1], 10, 64); err != nil {
		return clientErrorResponse("idTxc debe ser un número " + args[1])
	}
	var tipoTxc int
	if tipoTxc, err = strconv.Atoi(args[2]); err != nil || (tipoTxc != 1 && tipoTxc != 2) {
		return clientErrorResponse("tipoTxc debe ser un número con valor (1 ó 2)" + args[2])
	}

	var idOrg int
	if idOrg, err = strconv.Atoi(args[3]); err != nil {
		return clientErrorResponse("idOrg debe ser un número " + args[3])
	}

	var fechahoraRespuesta time.Time
	if fechahoraRespuesta, err = time.Parse(time.RFC3339, args[4]); err != nil {
		return clientErrorResponse("fechahoraRespuesta debe ser una fecha con formato yyyy-MM-dd'T'HH:mm:ss.SSSX " + args[4])
	}

	var tipoRespuesta int
	if tipoRespuesta, err = strconv.Atoi(args[5]); err != nil || (tipoRespuesta != TIPO_RESPUESTA_CONFIRMADO && tipoRespuesta != TIPO_RESPUESTA_RECHAZADO) {
		return clientErrorResponse("tipoRespuesta debe ser un número con valor (1 ó 2) " + args[5])
	}

	txc, _, errT := findTXConfirmable(APIstub, idOrg, uint64(idTxc))
	if errT != nil {
		return clientErrorResponse("No existe una TXConfirmable con id " + args[1] + " y organismo " + args[3])
	}

	if txc.IDOrganismo != idOrg {
		return clientErrorResponse("El idOrganismo no se corresponge con el organismo de la TXConfirmacion guardada")
	}

	if txc.TipoTxc != tipoTxc {
		return clientErrorResponse("El tipoTXC no se corresponge con el organismo de la TXConfirmacion guardada")
	}

	txc.FechaHoraRespuesta = fechahoraRespuesta
	txc.TipoRespuesta = tipoRespuesta

	// Creo un impuesto a partir de txc.assetValue (json) por ahora es solo Impuesto.
	var txc_pi Impuesto
	if err = json.Unmarshal([]byte(txc.AssetValue), &txc_pi); err != nil {
		return systemErrorResponse("Error unmarshal de txc.assetValue")
	}

	asset_pi, _, errAsset := findImpuesto(APIstub, txc.CUIT, txc_pi.Impuesto)
	existAsset := errAsset == nil
	if existAsset {
		if asset_pi.IDTxc != txc.IDTxc {
			return clientErrorResponse("El asset [v_key_pi] sobre el que debe impactar la TXC [v_key_txc] esta pendiente se ser actualizado por otra TXC [v_asset_pi.keyTxc]")
		}

		asset_pi.IDTxc = 0
		if txc.TipoRespuesta == TIPO_RESPUESTA_CONFIRMADO {
			asset_pi.Estado = txc_pi.Estado
			asset_pi.Periodo = txc_pi.Periodo
			asset_piAsBytes, _ := json.Marshal(asset_pi)
			log.Println("Guardando " + string(asset_piAsBytes))
			if err = APIstub.PutState(getImpuestoKeyByCuitId(txc.CUIT, asset_pi.Impuesto), asset_piAsBytes); err != nil {
				return systemErrorResponse("Error al guardar Impuesto - " + getImpuestoKeyByCuitId(txc.CUIT, asset_pi.Impuesto) + ", error " + err.Error())
			}

		}
	} else if txc.TipoRespuesta == TIPO_RESPUESTA_CONFIRMADO {

		if txc.TipoTxc == 1 {
			return clientErrorResponse("El asset [v_key_pi] que debe ser actualizado por la TXC [v_key_txc] no existe")
		}
		txc_piAsBytes, _ := json.Marshal(txc_pi)
		if err = APIstub.PutState(getImpuestoKeyByCuitId(txc.CUIT, txc_pi.Impuesto), txc_piAsBytes); err != nil {
			return systemErrorResponse("Error al guardar Impuesto - " + getImpuestoKeyByCuitId(txc.CUIT, txc_pi.Impuesto) + ", error " + err.Error())
		}
	}
	txcAsBytes, _ := json.Marshal(txc)
	if err = APIstub.PutState(getTxConfirmableKey(txc.IDOrganismo, txc.IDTxc), txcAsBytes); err != nil {
		return systemErrorResponse("Error al guardar TXConfirmable - " + getTxConfirmableKey(txc.IDOrganismo, txc.IDTxc) + ", error " + err.Error())
	}
	return successResponse("", 0)
}
