package organizaciones

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
)

type Org struct {
	ID        uint64            `json:"id,omitempty"        csv:"ID"`
	MSPID     string            `json:"mspid,omitempty"     csv:"MSPID"`
	Nombre    string            `json:"nombre,omitempty"    csv:"NOMBRE"`
	CUIT      uint64            `json:"cuit,omitempty"      csv:"CUIT,omitempty"`
	Provincia *common.Provincia `json:"provincia,omitempty" csv:"PROVINCIA,omitempty"`
}
