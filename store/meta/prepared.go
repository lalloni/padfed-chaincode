package meta

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
)

const witnessTag = "wit"

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
		singleton := singleton
		err := prepareSingleton(&singleton, members, value)
		if err != nil {
			return nil, err
		}
		singletons[singleton.Tag] = &singleton
	}
	collections := map[string]*Collection{}
	for _, collection := range com.Collections {
		collection := collection
		err := prepareCollection(&collection, members, valueType)
		if err != nil {
			return nil, err
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
	if com.Copier == nil {
		com.Copier = reflectionShallowCopy
	}
	return &PreparedComposite{
		name:        com.Name,
		composite:   &com,
		singletons:  singletons,
		collections: collections,
	}, nil
}

func prepareCollection(collection *Collection, members map[string]interface{}, valueType reflect.Type) error {
	if collection.Tag == "" {
		return errors.Errorf("composite collection %+v must specifify a tag name", collection)
	}
	if collection.Tag == witnessTag {
		return errors.Errorf("reserved member tag: collection %+v", collection)
	}
	if _, ok := members[collection.Tag]; ok {
		return errors.Errorf("duplicate member tag: collection %+v", collection)
	}
	if collection.Collector == nil {
		if collection.Field != "" {
			collection.Collector = MapCollector(FieldGetter(collection.Field), FieldSetter(collection.Field))
		} else {
			return errors.Errorf("composite collection with tag %q must have a collector function or specify a field name", collection.Tag)
		}
	}
	if collection.Enumerator == nil {
		if collection.Field != "" {
			collection.Enumerator = MapEnumerator(FieldGetter(collection.Field))
		} else {
			return errors.Errorf("composite collection with tag %q must have an enumerator function or specify a field name", collection.Tag)
		}
	}
	if collection.Creator == nil {
		if collection.Field != "" {
			field, ok := valueType.FieldByName(collection.Field)
			if !ok {
				return errors.Errorf("composite collection with tag %q field name %q does not match any value field", collection.Tag, collection.Field)
			}
			collection.Creator = ValueCreator(reflect.New(field.Type.Elem()).Elem().Interface())
		} else {
			return errors.Errorf("composite collection with tag %q must have a creator function or specify a field name", collection.Tag)
		}
	}
	if collection.Clear == nil {
		if collection.Field != "" {
			collection.Clear = FieldClear(collection.Field)
		} else {
			return errors.Errorf("composite collection with tag %q must have a clear function or specify a field name", collection.Tag)
		}
	}
	return nil
}

func prepareSingleton(singleton *Singleton, members map[string]interface{}, value interface{}) error {
	if singleton.Tag == "" {
		return errors.Errorf("composite singleton %+v must specifify a tag name", singleton)
	}
	if singleton.Tag == witnessTag {
		return errors.Errorf("reserved member tag: singleton %+v", singleton)
	}
	if _, ok := members[singleton.Tag]; ok {
		return errors.Errorf("duplicate member tag: singleton %+v", singleton)
	}
	if singleton.Getter == nil {
		if singleton.Field != "" {
			singleton.Getter = FieldGetter(singleton.Field)
		} else {
			return errors.Errorf("composite singleton with tag %q must have a getter function or specify a field name", singleton.Tag)
		}
	}
	if singleton.Setter == nil {
		if singleton.Field != "" {
			singleton.Setter = FieldSetter(singleton.Field)
		} else {
			return errors.Errorf("composite singleton with tag %q must have a setter function or specify a field name", singleton.Tag)
		}
	}
	if singleton.Creator == nil {
		if singleton.Field != "" {
			singleton.Creator = ValueCreator(FieldGetter(singleton.Field)(value))
		} else {
			return errors.Errorf("composite singleton with tag %q must have a creator function or specify a field name", singleton.Tag)
		}
	}
	if singleton.Clear == nil {
		if singleton.Field != "" {
			singleton.Clear = FieldClear(singleton.Field)
		} else {
			return errors.Errorf("composite singleton with tag %q must have a clear function or specify a field name", singleton.Tag)
		}
	}
	return nil
}

type PreparedComposite struct {
	name        string
	composite   *Composite
	singletons  map[string]*Singleton
	collections map[string]*Collection
}

func (cc *PreparedComposite) Name() string {
	return cc.name
}

func (cc *PreparedComposite) IdentifierKey(id interface{}) (k *key.Key, err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("building composite %q key from id %v: %v", cc.name, id, p)
		}
	}()
	return cc.composite.IdentifierKey(id)
}

func (cc *PreparedComposite) KeyIdentifier(k *key.Key) (v interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("building composite %q id from key %s: %v", cc.name, k, p)
		}
	}()
	return cc.composite.KeyIdentifier(k)
}

func (cc *PreparedComposite) ValueKey(val interface{}) (*key.Key, error) {
	id, err := cc.ValueIdentifier(val)
	if err != nil {
		return nil, err
	}
	return cc.IdentifierKey(id)
}

func (cc *PreparedComposite) ValueIdentifier(val interface{}) (id interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("getting composite %q id: %v", cc.name, p)
		}
	}()
	id = cc.composite.IdentifierGetter(val)
	return
}

func (cc *PreparedComposite) ValueWitness(val interface{}) (*Entry, error) {
	k, err := cc.ValueKey(val)
	if err != nil {
		return nil, err
	}
	return &Entry{
		Key:   k.Tagged(witnessTag),
		Value: 1,
	}, nil
}

func (cc *PreparedComposite) KeyWitness(key *key.Key) *key.Key {
	return key.Tagged(witnessTag)
}

func (cc *PreparedComposite) IsWitnessKey(key *key.Key) bool {
	return key.Tag.Name == witnessTag
}

func (cc *PreparedComposite) MustKeepRoot(val interface{}) bool {
	return cc.composite.KeepRoot
}

func (cc *PreparedComposite) RootEntry(val interface{}) (entry *Entry, err error) {
	valkey, err := cc.ValueKey(val)
	if err != nil {
		return nil, err
	}
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("getting composite %q root value: %v", cc.name, p)
		}
	}()
	root := cc.Cleared(val)
	entry = &Entry{Key: valkey, Value: root}
	return
}

func (cc *PreparedComposite) SingletonsEntries(val interface{}) (entries []*Entry, err error) {
	valkey, err := cc.ValueKey(val)
	if err != nil {
		return nil, err
	}
	var singleton *Singleton
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("getting composite %q singleton %q value: %v", cc.name, singleton.Tag, p)
		}
	}()
	entries = []*Entry{}
	for _, singleton = range cc.singletons {
		entries = append(entries, &Entry{
			Key:   valkey.Tagged(singleton.Tag),
			Value: singleton.Getter(val),
		})
	}
	return
}

func (cc *PreparedComposite) CollectionsEntries(val interface{}) (entries []*Entry, err error) {
	valkey, err := cc.ValueKey(val)
	if err != nil {
		return nil, err
	}
	var collection *Collection
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("getting composite %q collection %q items: %v", cc.name, collection.Tag, p)
		}
	}()
	entries = []*Entry{}
	for _, collection = range cc.collections {
		items := collection.Enumerator(val)
		for _, item := range items {
			entries = append(entries, &Entry{
				Key:   valkey.Tagged(collection.Tag, item.Identifier),
				Value: item.Value,
			})
		}
	}
	return
}

func (cc *PreparedComposite) Cleared(v interface{}) interface{} {
	nv := cc.Copy(v)
	for _, singleton := range cc.singletons {
		singleton.Clear(nv)
	}
	for _, collection := range cc.collections {
		collection.Clear(nv)
	}
	return nv
}

func (cc *PreparedComposite) Copy(v interface{}) interface{} {
	return cc.composite.Copier(v)
}

func (cc *PreparedComposite) Create() (v interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("creating composite %q value: %v", cc.name, p)
		}
	}()
	v = cc.composite.Creator()
	return
}

func (cc *PreparedComposite) SetIdentifier(val, id interface{}) (err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = errors.Errorf("setting composite %q id %v: %v", cc.name, id, p)
		}
	}()
	if cc.composite.IdentifierSetter != nil {
		cc.composite.IdentifierSetter(val, id)
	}
	return
}

func (cc *PreparedComposite) Collection(k *key.Key) *Collection {
	return cc.collections[k.Tag.Name]
}

func (cc *PreparedComposite) Singleton(k *key.Key) *Singleton {
	return cc.singletons[k.Tag.Name]
}

func (cc *PreparedComposite) KeyBaseName() string {
	return cc.composite.KeyBaseName
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

func reflectionShallowCopy(src interface{}) interface{} {
	ptr := false
	sv := reflect.ValueOf(src)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
		ptr = true
	}
	st := sv.Type()
	nv := reflect.New(st)
	if nv.Kind() == reflect.Ptr {
		nv = nv.Elem()
	}
	for f := 0; f < st.NumField(); f++ {
		v := sv.Field(f)
		nv.Field(f).Set(v)
	}
	if ptr {
		return nv.Addr().Interface()
	}
	return nv.Interface()
}
