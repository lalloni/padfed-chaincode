package store

import (
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/filtering"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/marshaling"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

type Range struct {
	First interface{}
	Last  interface{}
}

type Store interface {
	PutValue(key *key.Key, val interface{}) error
	GetValue(key *key.Key, val interface{}) (bool, error)
	HasValue(key *key.Key) (bool, error)
	DelValue(key *key.Key) error

	PutComposite(com *meta.PreparedComposite, val interface{}) error
	GetComposite(com *meta.PreparedComposite, id interface{}) (interface{}, error)
	HasComposite(com *meta.PreparedComposite, id interface{}) (bool, error)
	DelComposite(com *meta.PreparedComposite, id interface{}) error
	DelCompositeRange(com *meta.PreparedComposite, r Range) ([]interface{}, error)
}

func New(stub shim.ChaincodeStubInterface, opts ...Option) Store {
	s := &simplestore{
		stub:       stub,
		marshaling: DefaultMarshaling,
		filtering:  DefaultFiltering,
		sep:        key.DefaultSep,
		log:        shim.NewLogger("store"),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type simplestore struct {
	stub       shim.ChaincodeStubInterface
	log        *shim.ChaincodeLogger
	marshaling marshaling.Marshaling
	filtering  filtering.Filtering
	sep        *key.Sep
	seterrs    bool
}

func (s *simplestore) PutValue(k *key.Key, value interface{}) error {
	return s.internalPutValue(k, value)
}

func (s *simplestore) GetValue(k *key.Key, value interface{}) (bool, error) {
	return s.internalGetValue(k, value)
}

func (s *simplestore) HasValue(k *key.Key) (bool, error) {
	return s.internalHasValue(k)
}

func (s *simplestore) DelValue(k *key.Key) error {
	return s.internalDelValue(k)
}

func (s *simplestore) PutComposite(com *meta.PreparedComposite, val interface{}) error {
	we := com.ValueWitness(val)

	exist, err := s.internalHasValue(we.Key)
	if err != nil {
		return errors.Wrapf(err, "checking composite %q witness existence", com.Name)
	}
	if !exist {
		if err := s.internalPutValue(we.Key, we.Value); err != nil {
			return errors.Wrapf(err, "putting composite %q witness", com.Name)
		}
	}

	for _, entry := range com.SingletonsEntries(val) {
		if !reflect.ValueOf(entry.Value).IsNil() {
			if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q singleton %q", com.Name, entry)
			}
		}
	}
	for _, entry := range com.CollectionsEntries(val) {
		if reflect.ValueOf(entry.Value).IsNil() {
			if err := s.internalDelValue(entry.Key); err != nil {
				return errors.Wrapf(err, "deleting composite %q collection entry %q", com.Name, entry)
			}
		} else {
			if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q collection entry %q", com.Name, entry)
			}
		}
	}
	return nil
}

func (s *simplestore) GetComposite(com *meta.PreparedComposite, id interface{}) (interface{}, error) {
	valkey := com.IdentifierKey(id)
	if ok, err := s.HasComposite(com, id); err != nil {
		return nil, errors.Wrapf(err, "checking composite %q with key %q existence", com.Name, valkey)
	} else if !ok {
		return nil, nil // no existe la persona
	}
	val := com.Create()
	com.SetIdentifier(val, id)
	states, err := s.stub.GetStateByRange(valkey.RangeUsing(s.sep))
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q with key %q states iterator", com.Name, valkey)
	}
	merrs := []meta.MemberError{}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q with key %q next state", com.Name, valkey)
		}
		statekey, err := key.ParseUsing(state.GetKey(), s.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing composite %q with key %q item", com.Name, state.GetKey())
		}
		switch {
		case com.Collection(statekey) != nil:
			member := com.Collection(statekey)
			itemval := member.Creator()
			err := s.internalParseValue(state.GetValue(), itemval)
			if err != nil {
				s.log.Errorf("parsing composite %q with key %q collection item %q value in tx %s: %v", com.Name, valkey, statekey, s.stub.GetTxID(), err)
				if s.seterrs {
					seterr(itemval, err)
				}
				merrs = append(merrs, meta.MemberError{
					Kind:  "collection",
					Tag:   member.Tag,
					ID:    statekey.Tag.Value,
					Error: err.Error(),
				})
			}
			member.Collector(val, meta.Item{Identifier: statekey.Tag.Value, Value: itemval})
		case com.Singleton(statekey) != nil:
			member := com.Singleton(statekey)
			itemval := member.Creator()
			err := s.internalParseValue(state.GetValue(), itemval)
			if err != nil {
				s.log.Errorf("parsing composite %q with key %q collection item %q value in tx %s: %v", com.Name, valkey, statekey, s.stub.GetTxID(), err)
				if s.seterrs {
					seterr(itemval, err)
				}
				merrs = append(merrs, meta.MemberError{
					Kind:  "singleton",
					Tag:   member.Tag,
					Error: err.Error(),
				})
			}
			member.Setter(val, itemval)
		}
	}
	if s.seterrs && len(merrs) > 0 {
		seterrs(val, merrs)
	}
	return val, nil
}

func (s *simplestore) HasComposite(com *meta.PreparedComposite, id interface{}) (bool, error) {
	key := com.IdentifierKey(id)
	wk := com.KeyWitness(key)
	var a interface{}
	found, err := s.internalGetValue(wk, &a)
	if err != nil {
		return false, errors.Wrapf(err, "getting composite %q witness with key %q", com.Name, wk.StringUsing(s.sep))
	}
	return found, nil
}

func (s *simplestore) DelComposite(com *meta.PreparedComposite, id interface{}) error {
	key := com.IdentifierKey(id)
	states, err := s.stub.GetStateByRange(key.RangeUsing(s.sep))
	if err != nil {
		return errors.Wrapf(err, "getting composite %q states with key %q for deletion", com.Name, key)
	}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return errors.Wrapf(err, "getting composite %q with key %q next state for deletion", com.Name, key)
		}
		err = s.stub.DelState(state.GetKey())
		if err != nil {
			return errors.Wrapf(err, "deleting composite %q with key %q state %q", com.Name, key, state.GetKey())
		}
	}
	return nil
}

func (s *simplestore) DelCompositeRange(com *meta.PreparedComposite, r Range) ([]interface{}, error) {
	first, _ := com.IdentifierKey(r.First).RangeUsing(s.sep)
	_, last := com.IdentifierKey(r.Last).RangeUsing(s.sep)
	states, err := s.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] for deletion", com.Name, first, last)
	}
	defer states.Close()
	res := []interface{}{}
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] next key for deletion", com.Name, first, last)
		}
		statekey, err := key.ParseUsing(state.GetKey(), s.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing state key %q as composite %q key", state.GetKey(), com.Name)
		}
		if com.IsWitnessKey(statekey) {
			id, err := com.KeyIdentifier(statekey)
			if err != nil {
				return nil, errors.Wrapf(err, "getting composite %q identifier for state %q", com.Name, state.GetKey())
			}
			res = append(res, id)
		}
		err = s.stub.DelState(state.GetKey())
		if err != nil {
			return nil, errors.Wrapf(err, "deleting composite %q range [%q,%q] state %q", com.Name, first, last, state.GetKey())
		}
	}
	return res, nil
}

func (s *simplestore) internalPutValue(k *key.Key, value interface{}) error {
	if err := k.ValidateUsing(s.sep); err != nil {
		return errors.Wrap(err, "checking value key")
	} else if bs, err := s.marshaling.Marshal(value); err != nil {
		return errors.Wrap(err, "marshaling value")
	} else if bs, err := s.filtering.Filter(bs); err != nil {
		return errors.Wrap(err, "filtering value")
	} else if err := s.stub.PutState(k.StringUsing(s.sep), bs); err != nil {
		return errors.Wrap(err, "putting marshaled value into state")
	}
	return nil
}

func (s *simplestore) internalHasValue(k *key.Key) (bool, error) {
	bs, err := s.stub.GetState(k.StringUsing(s.sep))
	if err != nil {
		return false, errors.Wrap(err, "getting value from state")
	}
	return bs != nil, nil
}

func (s *simplestore) internalGetValue(k *key.Key, value interface{}) (bool, error) {
	bs, err := s.stub.GetState(k.StringUsing(s.sep))
	if err != nil {
		return false, errors.Wrap(err, "getting marshaled value from state")
	}
	if bs == nil {
		return false, nil
	}
	return true, s.internalParseValue(bs, value)
}

func (s *simplestore) internalDelValue(k *key.Key) error {
	err := s.stub.DelState(k.StringUsing(s.sep))
	if err != nil {
		return errors.Wrap(err, "deleting value from state")
	}
	return nil
}

func (s *simplestore) internalParseValue(bs []byte, value interface{}) error {
	if bs, err := s.filtering.Unfilter(bs); err != nil {
		return errors.Wrap(err, "unfiltering value")
	} else if err := s.marshaling.Unmarshal(bs, value); err != nil {
		return errors.Wrap(err, "unmarshaling value")
	}
	return nil
}

func seterr(val interface{}, e error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Error")
		if f.Kind() == reflect.String {
			f.SetString(e.Error())
		}
	}
}

func seterrs(val interface{}, merrs []meta.MemberError) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Errors")
		if f.Kind() == reflect.Interface {
			f.Set(reflect.ValueOf(merrs))
		}
	}
}
