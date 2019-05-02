package meta

import (
	"fmt"
	"reflect"
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
	valueType := reflect.TypeOf(value).Elem()
	members := map[string]interface{}{}
	singletons := map[string]*Singleton{}
	for _, singleton := range com.Singletons {
		singleton := singleton // clone for mutation
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
		singletons[singleton.Tag] = &singleton
	}
	collections := map[string]*Collection{}
	for _, collection := range com.Collections {
		collection := collection // clone for mutation
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
				field, ok := valueType.FieldByName(collection.Field)
				if !ok {
					return nil, errors.Errorf("composite collection with tag %q field name %q does not match any value field", collection.Tag, collection.Field)
				}
				collection.Creator = ValueCreator(reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				return nil, errors.Errorf("composite collection with tag %q must have a creator function or specify a field name", collection.Tag)
			}
		}
		collections[collection.Tag] = &collection
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
		Name:        com.Name,
		composite:   &com,
		singletons:  singletons,
		collections: collections,
	}, nil
}

type PreparedComposite struct {
	Name        string
	composite   *Composite
	singletons  map[string]*Singleton
	collections map[string]*Collection
}

func (cc *PreparedComposite) IdentifierKey(id interface{}) *key.Key {
	return cc.composite.Keyer(id)
}

func (cc *PreparedComposite) ValueKey(val interface{}) *key.Key {
	return cc.IdentifierKey(cc.composite.IdentifierGetter(val))
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
	for _, singleton := range cc.singletons {
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
	for _, collection := range cc.collections {
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

func (cc *PreparedComposite) Range(key *key.Key) (string, string) {
	s := key.String()
	return s, s + string(utf8.MaxRune)
}

func (cc *PreparedComposite) RangeSep(key *key.Key, sep *key.Sep) (string, string) {
	s := key.StringUsing(sep)
	return s, s + string(utf8.MaxRune)
}

func (cc *PreparedComposite) Collection(key *key.Key) *Collection {
	return cc.collections[key.Tag.Name]
}

func (cc *PreparedComposite) Singleton(key *key.Key) *Singleton {
	return cc.singletons[key.Tag.Name]
}

func (cc *PreparedComposite) SetIdentifier(val, id interface{}) {
	if cc.composite.IdentifierSetter != nil {
		cc.composite.IdentifierSetter(val, id)
	}
}

func (cc *PreparedComposite) CreateValue() interface{} {
	return cc.composite.Creator()
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
