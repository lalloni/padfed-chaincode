package response

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

type Response struct {
	Status  int32
	Message string
	Payload *Payload
}

type Payload struct {
	Client      *Client      `json:"client,omitempty"`
	Transaction *Transaction `json:"transation,omitempty"`
	Result      interface{}  `json:"result,omitempty"`
	Fault       interface{}  `json:"fault,omitempty"`
}

type Client struct {
	MSPID   string `json:"mspid,omitempty"`
	Subject string `json:"subject,omitempty"`
	Issuer  string `json:"issuer,omitempty"`
}

type Transaction struct {
	ID       string           `json:"id,omitempty"`
	Function context.Function `json:"function,omitempty"`
}

func (r *Response) OK() bool {
	return r.Status < shim.ERRORTHRESHOLD
}
