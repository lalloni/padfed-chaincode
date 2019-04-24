package meta

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
)

type KeyerFunc func(interface{}) *key.Key

type CreatorFunc func() interface{}

type GetterFunc func(interface{}) interface{}
type SetterFunc func(interface{}, interface{})

type EnumeratorFunc func(interface{}) []Item
type CollectorFunc func(interface{}, Item)

type Item struct {
	Identifier string
	Value      interface{}
}

type Composite struct {
	Name            string
	Creator         CreatorFunc
	Identifier      GetterFunc
	IdentifierField string
	Keyer           KeyerFunc
	Singletons      []Singleton
	Collections     []Collection
}

type Singleton struct {
	Tag     string
	Field   string
	Creator CreatorFunc
	Getter  GetterFunc
	Setter  SetterFunc
}

type Collection struct {
	Tag        string
	Field      string
	Creator    CreatorFunc
	Collector  CollectorFunc
	Enumerator EnumeratorFunc
}
