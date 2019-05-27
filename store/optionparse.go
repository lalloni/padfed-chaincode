package store

import (
	"net/url"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
)

func Options(stub shim.ChaincodeStubInterface) ([]Option, error) {
	ss := strings.SplitN(string(stub.GetArgs()[0]), "?", 2)
	if len(ss) < 2 {
		return nil, nil
	}
	q, err := url.ParseQuery(ss[1])
	if err != nil {
		return nil, errors.Wrap(err, "parsing function options")
	}
	oo := []Option{}
	for k := range q {
		switch k {
		case "lenientread":
			oo = append(oo, SetLenient(true))
		case "embederrors":
			oo = append(oo, SetErrors(true))
		}
	}
	return oo, nil
}
