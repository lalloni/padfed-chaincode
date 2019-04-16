package meta

import (
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
)

var witnessTag = "wit"

func MustPrepare(com Composite) *PreparedComposite {
	cc, err := Prepare(com)
	if err != nil {
		panic(err)
	}
	return cc
}

func Prepare(com Composite) (*PreparedComposite, error) {
	members := map[string]interface{}{}
	for _, singleton := range com.Singletons {
		if singleton.Tag == witnessTag {
			return nil, errors.Errorf("reserved member tag: singleton %+v", singleton)
		}
		if _, ok := members[singleton.Tag]; ok {
			return nil, errors.Errorf("duplicate member tag: singleton %+v", singleton)
		}
		members[singleton.Tag] = singleton
	}
	for _, collection := range com.Collections {
		if collection.Tag == witnessTag {
			return nil, errors.Errorf("reserved member tag: collection %+v", collection)
		}
		if _, ok := members[collection.Tag]; ok {
			return nil, errors.Errorf("duplicate member tag: collection %+v", collection)
		}
		members[collection.Tag] = collection
	}
	return &PreparedComposite{
		Name:      com.Name,
		Composite: &com,
		Members:   members,
	}, nil
}

type PreparedComposite struct {
	Name      string
	Composite *Composite
	Members   map[string]interface{}
}

func (cc *PreparedComposite) IdentifierKey(id interface{}) *key.Key {
	return cc.Composite.Keyer(id)
}

func (cc *PreparedComposite) ValueKey(val interface{}) *key.Key {
	return cc.IdentifierKey(cc.Composite.Identifier(val))
}

func (cc *PreparedComposite) ValueWitness(val interface{}) *Entry {
	return &Entry{
		Key:   cc.ValueKey(val).Tagged(witnessTag),
		Value: 1,
	}
}

func (cc *PreparedComposite) KeyWitness(key *key.Key) *key.Key {
	return key.Tagged(witnessTag)
}

func (cc *PreparedComposite) SingletonsEntries(val interface{}) []*Entry {
	valkey := cc.ValueKey(val)
	entries := []*Entry(nil)
	for _, singleton := range cc.Composite.Singletons {
		entries = append(entries, &Entry{
			Key:   valkey.Tagged(singleton.Tag),
			Value: singleton.Getter(val),
		})
	}
	return entries
}

func (cc *PreparedComposite) CollectionsEntries(val interface{}) []*Entry {
	valkey := cc.ValueKey(val)
	entries := []*Entry(nil)
	for _, collection := range cc.Composite.Collections {
		items := []Item{}
		collection.Enumerator(val, &items)
		for _, item := range items {
			entries = append(entries, &Entry{
				Key:   valkey.Tagged(collection.Tag, item.Identifier),
				Value: item.Value,
			})
		}
	}
	return entries
}

func (cc *PreparedComposite) Range(key *key.Key, sep *key.Sep) (string, string) {
	s := key.StringUsing(sep)
	return s, s + string(utf8.MaxRune)
}

func (cc *PreparedComposite) Member(key *key.Key) interface{} {
	return cc.Members[key.Tag.Name]
}

type Entry struct {
	Key   *key.Key
	Value interface{}
}

func (e *Entry) String() string {
	return fmt.Sprintf("[%s → %+v]", e.Key, e.Value)
}

func NewItem(id string, value interface{}) Item {
	return Item{
		Identifier: id,
		Value:      value,
	}
}
