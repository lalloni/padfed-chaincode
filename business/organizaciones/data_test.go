package organizaciones

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
)

func TestLoad(t *testing.T) {
	a := assert.New(t)
	orgs, err := load()
	a.NoError(err)
	a.EqualValues(27, len(orgs))
	first := []Org{
		{ID: 1, MSPID: "AFIP", Nombre: "AFIP", CUIT: 33693450239, Provincia: nil},
		{ID: 100, MSPID: "MORGS", Nombre: "MULTI ORGANISMOS", CUIT: 0, Provincia: nil},
		{ID: 900, MSPID: "COMARB", Nombre: "COMISIÓN ARBITRAL DEL CONVENIO MULTILATERAL", CUIT: 30658892718, Provincia: nil},
		{ID: 901, MSPID: "AGIP", Nombre: "AGIP - CABA", CUIT: 34999032089, Provincia: common.ProvinciaCódigo(0)},
		{ID: 902, MSPID: "ARBA", Nombre: "ARBA - BUENOS AIRES", CUIT: 30710404611, Provincia: common.ProvinciaCódigo(1)},
	}
	a.EqualValues(first, orgs[0:len(first)])
}
