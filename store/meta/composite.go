package meta

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
)

type KeyFunc func(interface{}) *key.Key
type CreateFunc func() interface{}
type GetFunc func(interface{}) interface{}
type SetFunc func(interface{}, interface{})
type EnumerateFunc func(interface{}, *[]Item)
type CollectFunc func(interface{}, Item)

type Item struct {
	Identifier string
	Value      interface{}
}

type Composite struct {
	Name        string
	Creator     CreateFunc
	Identifier  GetFunc
	Keyer       KeyFunc
	Singletons  []Singleton
	Collections []Collection
}

type Singleton struct {
	Tag     string
	Creator CreateFunc
	Getter  GetFunc
	Setter  SetFunc
}

type Collection struct {
	Tag        string
	Creator    CreateFunc
	Collector  CollectFunc
	Enumerator EnumerateFunc
}
