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
	value := com.Creator()
	members := map[string]interface{}{}
	for _, singleton := range com.Singletons {
		if singleton.Tag == witnessTag {
			return nil, errors.Errorf("reserved member tag: singleton %+v", singleton)
		}
		if _, ok := members[singleton.Tag]; ok {
			return nil, errors.Errorf("duplicate member tag: singleton %+v", singleton)
		}
		if singleton.Getter == nil {
			if singleton.Field != "" {
				singleton.Getter = FieldGetter(singleton.Field)
			} else {
				return nil, errors.Errorf("composite singleton with tag %q must have a getter function or specify a field name", singleton.Tag)
			}
		}
		if singleton.Setter == nil {
			if singleton.Field != "" {
				singleton.Setter = FieldSetter(singleton.Field)
			} else {
				return nil, errors.Errorf("composite singleton with tag %q must have a setter function or specify a field name", singleton.Tag)
			}
		}
		if singleton.Creator == nil {
			if singleton.Field != "" {
				singleton.Creator = ValueCreator(FieldGetter(singleton.Field)(value))
			} else {
				return nil, errors.Errorf("composite singleton with tag %q must have a creator function or specify a field name", singleton.Tag)
			}
		}
		members[singleton.Tag] = singleton
	}
	for _, collection := range com.Collections {
		if collection.Tag == "" {
			return nil, errors.Errorf("composite collection %+v must specifify a tag name", collection)

		}
		if collection.Tag == witnessTag {
			return nil, errors.Errorf("reserved member tag: collection %+v", collection)
		}
		if _, ok := members[collection.Tag]; ok {
			return nil, errors.Errorf("duplicate member tag: collection %+v", collection)
		}
		if collection.Collector == nil {
			if collection.Field != "" {
				collection.Collector = MapCollector(FieldGetter(collection.Field))
			} else {
				return nil, errors.Errorf("composite collection with tag %q must have a collector function or specify a field name", collection.Tag)
			}
		}
		if collection.Enumerator == nil {
			if collection.Field != "" {
				collection.Enumerator = MapEnumerator(FieldGetter(collection.Field))
			} else {
				return nil, errors.Errorf("composite collection with tag %q must have an enumerator function or specify a field name", collection.Tag)
			}
		}
		if collection.Creator == nil {
			if collection.Field != "" {
				collection.Creator = ValueCreator(FieldGetter(collection.Field)(value))
			} else {
				return nil, errors.Errorf("composite collection with tag %q must have a creator function or specify a field name", collection.Tag)
			}
		}
		members[collection.Tag] = collection
	}
	if com.IdentifierGetter == nil {
		if com.IdentifierField != "" {
			com.IdentifierGetter = FieldGetter(com.IdentifierField)
		} else {
			return nil, errors.New("composite must have an identifier getter function or specify an identifier field name")
		}
	}
	if com.IdentifierSetter == nil {
		if com.IdentifierField != "" {
			com.IdentifierSetter = FieldSetter(com.IdentifierField)
		}
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
	return cc.IdentifierKey(cc.Composite.IdentifierGetter(val))
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
		items := collection.Enumerator(val)
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
