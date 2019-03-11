package fabric

import (
	"bytes"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func QueryHistory(stub shim.ChaincodeStubInterface, args []string) *Response {
	if len(args) != 1 {
		return ClientErrorResponse("Número incorrecto de parámetros. Se esperaba 1 con {asset_key}")
	}
	pKey := args[0]
	resultsIterator, err := stub.GetHistoryForKey(pKey)
	if err != nil {
		return SystemErrorResponse(err.Error())
	}
	defer resultsIterator.Close()

	// TODO: eliminar construcción manual de JSON
	var buffer bytes.Buffer
	buffer.WriteString("[")
	count := 0
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return SystemErrorResponse(err.Error())
		}
		count++

		if bArrayMemberAlreadyWritten {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(pKey + "." + strconv.Itoa(count))
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	Log.Info("- getHistory returning:" + buffer.String())

	return SuccessResponseWithBuffer(&buffer)
}
