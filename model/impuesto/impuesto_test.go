package impuesto

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

	tx := test.MockTransactionStart(t, mock)
	err := st.PutComposite(Schema, imp1)
	a.NoError(err)
	test.MockTransactionEnd(t, mock, tx)

	tx = test.MockTransactionStart(t, mock)
	imp2, err := st.GetComposite(Schema, uint64(10))
	a.NoError(err)
	a.EqualValues(imp1, imp2)
	test.MockTransactionEnd(t, mock, tx)

}
