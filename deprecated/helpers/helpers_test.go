package helpers_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/deprecated/helpers"
)

func TestGetCUIT(t *testing.T) {
	_, err := helpers.GetCUIT("1")
	if err == nil {
		t.Error("Debe pinchar, el valor 1 no es un cuit valido.")
	}
	_, err = helpers.GetCUIT("20066675573")
	if err != nil {
		t.Errorf("No debe pinchar, el valor cuit es valido: %v", err)
	}
}
