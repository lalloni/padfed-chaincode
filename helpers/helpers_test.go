package helpers_test

import (
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode/helpers"
)

func TestGetCUIT(t *testing.T) {
	_, err := helpers.GetCUIT("1")
	if err == nil {
		t.Error("Debe pinchar, el valor 1 no es un cuit valido.")
	}
	_, err = helpers.GetCUIT("20066675573")
	if err != nil {
		t.Error("No debe pinchar, el valor cuit es valido.")
	}
}
