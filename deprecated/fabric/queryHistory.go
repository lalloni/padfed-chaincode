package fabric

import (
	"bytes"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
)

func QueryHistory(stub shim.ChaincodeStubInterface, _ []string) *Response {
	args, err := handler.ExtractArgs(stub.GetArgs()[1:], param.String)
	if err != nil {
		return ClientErrorResponse(err.Error())
	}

	pKey := args[0].(string)
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
		buffer.WriteString(pKey + " [v" + strconv.Itoa(count) + "]")
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
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).In(time.Local).Format(time.RFC3339Nano))
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
