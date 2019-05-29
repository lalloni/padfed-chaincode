package fabric

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var per11d = regexp.MustCompile(`(^PER_)(\d{11})`)

func DeleteAll(stub shim.ChaincodeStubInterface) *Response {
	return DeleteByKeyRange(stub, "", "")
}

func DeleteByKeyRange(stub shim.ChaincodeStubInterface, startKey, endKey string) *Response {
	var firstDeleted, lastDeleted, partialMsj string
	log.Println("Se eliminaran los primeros 100 elementos para evitar un timeout")
	total := 100
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return SystemErrorResponse(err.Error())
	}
	count := 0
	for resultsIterator.HasNext() {
		count++
		if count > total {
			break
		}
		result, err := resultsIterator.Next()
		if err != nil {
			return SystemErrorResponse(err.Error(), count)
		}
		if firstDeleted != "" {
			lastDeleted = result.Key
		} else {
			firstDeleted = result.Key
			if !resultsIterator.HasNext() {
				lastDeleted = firstDeleted
			}
		}
		// lenKeyPER = "PER_" + CUIT = 15 digits
		const lenKeyPER = 15
		if len(result.Key) == lenKeyPER {
			// delete asset CUIT
			res := per11d.MatchString(result.Key)
			if res {
				splitedKey := strings.Split(result.Key, "_")
				err := stub.DelState(splitedKey[1])
				if err != nil {
					return SystemErrorResponse(err.Error())
				}
			}
		}
		log.Print("[" + strconv.Itoa(count) + "] key to delete [" + result.Key)
		err = stub.DelState(result.Key)
		if err != nil {
			return SystemErrorResponse(err.Error())
		}
	}
	if firstDeleted != "" {
		partialMsj = "Keys eliminadas desde [" + firstDeleted + "] hasta [" + lastDeleted + "]"
	}
	return SuccessResponse(partialMsj+" hasNext ["+strconv.FormatBool(count > total)+"]", count)
}
