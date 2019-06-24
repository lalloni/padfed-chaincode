package meta

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
)

type ValFunc func(*key.Key) (interface{}, error)
type KeyFunc func(interface{}) (*key.Key, error)

type CreatorFunc func() (new interface{})
type CopierFunc func(src interface{}) (new interface{})

type GetterFunc func(src interface{}) interface{}
type SetterFunc func(tgt interface{}, v interface{})
type MutatorFunc func(tgt interface{})

type EnumeratorFunc func(src interface{}) []Item
type CollectorFunc func(tgt interface{}, i Item)

type Item struct {
	Identifier string
	Value      interface{}
}

type Composite struct {
	Name             string
	Creator          CreatorFunc
	Copier           CopierFunc
	IdentifierField  string
	IdentifierGetter GetterFunc
	IdentifierSetter SetterFunc
	IdentifierKey    KeyFunc
	KeyIdentifier    ValFunc
	KeyBaseName      string
	Singletons       []Singleton
	Collections      []Collection
	KeepRoot         bool
}

type Singleton struct {
	Tag     string
	Field   string
	Creator CreatorFunc
	Getter  GetterFunc
	Setter  SetterFunc
	Clear   MutatorFunc
}

type Collection struct {
	Tag        string
	Field      string
	Creator    CreatorFunc
	Collector  CollectorFunc
	Enumerator EnumeratorFunc
	Clear      MutatorFunc
}
