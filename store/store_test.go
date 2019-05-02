package store_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

type Thingy struct {
	Blah string `json:"blah,omitempty"`
}

type Thing struct {
	ID       uint64   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Age      uint8    `json:"age,omitempty"`
	Thingies []Thingy `json:"thingies,omitempty"`
}

type Other struct {
	Name   string `json:"name,omitempty"`
	Number int    `json:"number,omitempty"`
}

type Foo struct {
	Some string `json:"some,omitempty"`
	Num  int    `json:"num,omitempty"`
}

type Item struct {
	Name     string  `json:"name,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
}

type Compo struct {
	Thing *Thing           `json:"thing,omitempty"`
	Other *Other           `json:"other,omitempty"`
	Items map[string]*Item `json:"items,omitempty"`
	Foos  map[string]*Foo  `json:"foos,omitempty"`
}

var cc = meta.MustPrepare(meta.Composite{
	Name:    "compo",
	Creator: func() interface{} { return &Compo{Items: map[string]*Item{}, Foos: map[string]*Foo{}} },
	IdentifierGetter: func(v interface{}) interface{} {
		return v.(*Compo).Thing.ID
	},
	Keyer: func(id interface{}) *key.Key {
		return key.Based("compo", strconv.FormatUint(id.(uint64), 10))
	},
	Singletons: []meta.Singleton{
		{
			Tag:     "thing",
			Creator: func() interface{} { return &Thing{} },
			Getter:  func(v interface{}) interface{} { return v.(*Compo).Thing },
			Setter:  func(v interface{}, w interface{}) { v.(*Compo).Thing = w.(*Thing) },
		},
		{
			Tag:   "other",
			Field: "Other",
		},
	},
	Collections: []meta.Collection{
		{
			Tag:       "item",
			Creator:   func() interface{} { return &Item{} },
			Collector: func(v interface{}, i meta.Item) { v.(*Compo).Items[i.Identifier] = i.Value.(*Item) },
			Enumerator: func(v interface{}) []meta.Item {
				items := []meta.Item{}
				for k, v := range v.(*Compo).Items {
					items = append(items, meta.Item{Identifier: k, Value: v})
				}
				return items
			},
		},
		{
			Tag:   "foos",
			Field: "Foos",
		},
	},
})

func TestPutAndGetValue(t *testing.T) {
	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")
	a := assert.New(t)

	stub := shim.NewMockStub("test", nil)

	st := store.New(stub)
	key := key.Based("thing", "1")

	stub.MockTransactionStart("x")

	t1 := &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}}
	err := st.PutValue(key, t1)
	a.NoError(err)

	stub.MockTransactionEnd("x")

	t2 := &Thing{1234, "AA", 100, []Thingy{}}
	_, err = st.GetValue(key, t2)
	a.NoError(err)
	a.Equal(t1, t2)
}

func TestPutAndGet(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}},
		Other: &Other{"TT", 2123},
		Items: map[string]*Item{
			"a": {Name: "Pedro", Quantity: 10.0},
			"b": {Name: "Pablo", Quantity: 20.0},
		},
		Foos: map[string]*Foo{
			"foo1": {Some: "bar", Num: 634},
			"foo2": {Some: "baz", Num: 634},
		},
	}

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %+v", mustMarshal(c1))

	c2, err := st.GetComposite(cc, c1.Thing.ID)
	a.NoError(err)
	t.Logf("get: %+v", mustMarshal(c2))
	a.Equal(c1, c2)
}

func TestPutAndDelete(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}},
		Items: map[string]*Item{
			"a": {
				Name:     "Pedro",
				Quantity: 10.0,
			},
			"b": {
				Name:     "Pablo",
				Quantity: 20.0,
			},
		},
	}

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %+v", mustMarshal(c1))

	has, err := st.HasComposite(cc, c1.Thing.ID)
	a.NoError(err)
	a.True(has)

	err = st.DelComposite(cc, c1.Thing.ID)
	a.NoError(err)

	has, err = st.HasComposite(cc, c1.Thing.ID)
	a.NoError(err)
	a.False(has)

}

func mustMarshal(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
