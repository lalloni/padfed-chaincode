package impuestos

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/test"
)

func TestPut(t *testing.T) {

	a := assert.New(t)
	shim.SetLoggingLevel(shim.LogDebug)
	mock := test.NewMock("test", nil)
	st := store.New(mock)

	imp1 := &Impuesto{
		Codigo:      10,
		Org:         1,
		Nombre:      "Impuesto al Valor Agregado",
		Abreviatura: "IVA",
	}

	test.InTransaction(mock, func(tx string) {
		err := st.PutComposite(Schema, imp1)
		a.NoError(err)
	})

	imp2, err := st.GetComposite(Schema, uint64(10))
	a.NoError(err)
	a.EqualValues(imp1, imp2)

}
