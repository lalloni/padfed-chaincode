package store

import (
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/filtering"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/marshaling"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

var log = shim.NewLogger("store")

type Store interface {
	PutValue(key *key.Key, val interface{}) error
	GetValue(key *key.Key, val interface{}) (bool, error)
	HasValue(key *key.Key) (bool, error)
	DelValue(key *key.Key) error

	PutComposite(com *meta.PreparedComposite, val interface{}) error
	GetComposite(com *meta.PreparedComposite, id interface{}) (interface{}, error)
	HasComposite(com *meta.PreparedComposite, id interface{}) (bool, error)
	DelComposite(com *meta.PreparedComposite, id interface{}) error
	DelCompositeRange(com *meta.PreparedComposite, r *Range) ([]interface{}, error)
	GetCompositeRange(com *meta.PreparedComposite, r *Range) ([]interface{}, error)
	GetCompositeAll(com *meta.PreparedComposite) ([]interface{}, error)
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
	we, err := com.ValueWitness(val)
	if err != nil {
		return errors.Wrapf(err, "getting composite %q value witness", com.Name())
	}
	exist, err := s.internalHasValue(we.Key)
	if err != nil {
		return errors.Wrapf(err, "checking composite %q witness existence", com.Name())
	}
	if !exist {
		if err := s.internalPutValue(we.Key, we.Value); err != nil {
			return errors.Wrapf(err, "putting composite %q witness", com.Name())
		}
	}
	hascomps := false
	entries, err := com.SingletonsEntries(val)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(entries) > 0 {
		hascomps = true
	}
	for _, entry := range entries {
		if !reflect.ValueOf(entry.Value).IsNil() {
			if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q singleton %q", com.Name(), entry)
			}
		}
	}
	entries, err = com.CollectionsEntries(val)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(entries) > 0 {
		hascomps = true
	}
	for _, entry := range entries {
		if reflect.ValueOf(entry.Value).IsNil() {
			if err := s.internalDelValue(entry.Key); err != nil {
				return errors.Wrapf(err, "deleting composite %q collection entry %q", com.Name(), entry)
			}
		} else {
			if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q collection entry %q", com.Name(), entry)
			}
		}
	}
	if !hascomps || com.MustKeepRoot(val) {
		entry, err := com.RootEntry(val)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
			return errors.Wrapf(err, "putting composite %q root entry %q", com.Name(), entry)
		}
	}
	return nil
}

func (s *simplestore) GetComposite(com *meta.PreparedComposite, id interface{}) (interface{}, error) {
	valkey, err := com.IdentifierKey(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ok, err := s.HasComposite(com, id); err != nil {
		return nil, errors.Wrapf(err, "checking composite %q with key %q existence", com.Name(), valkey)
	} else if !ok {
		return nil, nil // no existe la persona
	}
	val, err := com.Create()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = com.SetIdentifier(val, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	states, err := s.stub.GetStateByRange(valkey.RangeUsing(s.sep))
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q with key %q states iterator", com.Name(), valkey)
	}
	merrs := []meta.MemberError{}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q with key %q next state", com.Name(), valkey)
		}
		statekey, err := key.ParseUsing(state.GetKey(), s.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing composite %q with key %q item", com.Name(), state.GetKey())
		}
		merr := s.inject(com, statekey, state, valkey, val)
		if s.seterrs && merr != nil {
			merrs = append(merrs, *merr)
		}
	}
	if s.seterrs && len(merrs) > 0 {
		seterrs(val, merrs)
	}
	return val, nil
}

func (s *simplestore) HasComposite(com *meta.PreparedComposite, id interface{}) (bool, error) {
	key, err := com.IdentifierKey(id)
	if err != nil {
		return false, errors.WithStack(err)
	}
	wk := com.KeyWitness(key)
	var a interface{}
	found, err := s.internalGetValue(wk, &a)
	if err != nil {
		return false, errors.Wrapf(err, "getting composite %q witness with key %q", com.Name(), wk.StringUsing(s.sep))
	}
	return found, nil
}

func (s *simplestore) DelComposite(com *meta.PreparedComposite, id interface{}) error {
	key, err := com.IdentifierKey(id)
	if err != nil {
		return errors.WithStack(err)
	}
	states, err := s.stub.GetStateByRange(key.RangeUsing(s.sep))
	if err != nil {
		return errors.Wrapf(err, "getting composite %q states with key %q for deletion", com.Name(), key)
	}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return errors.Wrapf(err, "getting composite %q with key %q next state for deletion", com.Name(), key)
		}
		err = s.stub.DelState(state.GetKey())
		if err != nil {
			return errors.Wrapf(err, "deleting composite %q with key %q state %q", com.Name(), key, state.GetKey())
		}
	}
	return nil
}

func (s *simplestore) DelCompositeRange(com *meta.PreparedComposite, r *Range) ([]interface{}, error) {
	first, last, err := s.identifierKeyRange(com, r)
	if err != nil {
		return nil, errors.Wrapf(err, "getting keys range %v", r)
	}
	states, err := s.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] for deletion", com.Name(), first, last)
	}
	defer states.Close()
	res := []interface{}{}
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] next key for deletion", com.Name(), first, last)
		}
		statekey, err := key.ParseUsing(state.GetKey(), s.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing state key %q as composite %q key", state.GetKey(), com.Name())
		}
		if com.IsWitnessKey(statekey) {
			id, err := com.KeyIdentifier(statekey)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			res = append(res, id)
		}
		err = s.stub.DelState(state.GetKey())
		if err != nil {
			return nil, errors.Wrapf(err, "deleting composite %q range [%q,%q] state %q", com.Name(), first, last, state.GetKey())
		}
	}
	return res, nil
}

func (s *simplestore) GetCompositeRange(com *meta.PreparedComposite, r *Range) ([]interface{}, error) {
	first, last, err := s.identifierKeyRange(com, r)
	if err != nil {
		return nil, errors.Wrapf(err, "getting keys range %v", r)
	}
	states, err := s.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] for reading", com.Name(), first, last)
	}
	defer states.Close()
	return s.internalReadCompositeIterator(com, states)
}

func (s *simplestore) GetCompositeAll(com *meta.PreparedComposite) ([]interface{}, error) {
	kbn := com.KeyBaseName()
	if kbn == "" {
		return nil, errors.Errorf("getting composite %q all instances: keybasename is empty", com.Name())
	}
	first, last := key.NewBase(kbn, "").RangeUsing(s.sep)
	states, err := s.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q all instances for reading", com.Name())
	}
	defer states.Close()
	return s.internalReadCompositeIterator(com, states)
}

func (s *simplestore) internalReadCompositeIterator(com *meta.PreparedComposite, states shim.StateQueryIteratorInterface) ([]interface{}, error) {
	var (
		valkey  *key.Key
		val, id interface{}
	)
	merrs := []meta.MemberError{}
	res := []interface{}{}
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q iterator next key for reading", com.Name())
		}
		statekey, err := key.ParseUsing(state.GetKey(), s.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing state key %q as composite %q key", state.GetKey(), com.Name())
		}
		basekey := key.NewBaseKey(statekey)
		if valkey == nil || !valkey.Equal(basekey) {
			if s.seterrs && len(merrs) > 0 {
				seterrs(val, merrs)
				merrs = []meta.MemberError{}
			}
			valkey = basekey
			val, err = com.Create()
			if err != nil {
				return nil, errors.WithStack(err)
			}
			id, err = com.KeyIdentifier(valkey)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			err = com.SetIdentifier(val, id)
			if err != nil {
				return nil, errors.Wrapf(err, "setting composite %q id %v from key %v", com.Name(), id, valkey)
			}
			res = append(res, val)
		}
		merr := s.inject(com, statekey, state, valkey, val)
		if s.seterrs && merr != nil {
			merrs = append(merrs, *merr)
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
	} else {
		ks := k.StringUsing(s.sep)
		if log.IsEnabledFor(shim.LogDebug) {
			log.Debugf("putting key '%s' with value '%s'", ks, string(bs))
		}
		if err := s.stub.PutState(ks, bs); err != nil {
			return errors.Wrap(err, "putting marshaled value into state")
		}
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

func (s *simplestore) inject(com *meta.PreparedComposite, statekey *key.Key, state *queryresult.KV, valkey *key.Key, val interface{}) *meta.MemberError {
	var merr *meta.MemberError
	switch {
	case statekey.Equal(valkey):
		err := s.internalParseValue(state.GetValue(), val)
		if err != nil {
			s.log.Errorf("parsing composite %q with key root item %q value in tx %s: %v", com.Name, valkey, s.stub.GetTxID(), err)
			if s.seterrs {
				seterr(val, err)
			}
			merr = &meta.MemberError{
				Kind:  "root",
				Error: err.Error(),
			}
		}
	case com.Collection(statekey) != nil:
		member := com.Collection(statekey)
		itemval := member.Creator()
		err := s.internalParseValue(state.GetValue(), itemval)
		if err != nil {
			s.log.Errorf("parsing composite %q with key %q collection item %q value in tx %s: %v", com.Name, valkey, statekey, s.stub.GetTxID(), err)
			if s.seterrs {
				seterr(itemval, err)
			}
			merr = &meta.MemberError{
				Kind:  "collection",
				Tag:   member.Tag,
				ID:    statekey.Tag.Value,
				Error: err.Error(),
			}
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
			merr = &meta.MemberError{
				Kind:  "singleton",
				Tag:   member.Tag,
				Error: err.Error(),
			}
		}
		member.Setter(val, itemval)
	}
	return merr
}

func (s *simplestore) identifierKeyRange(com *meta.PreparedComposite, r *Range) (string, string, error) {
	fk, err := com.IdentifierKey(r.First)
	if err != nil {
		return "", "", errors.Wrap(err, "getting range start key")
	}
	lk, err := com.IdentifierKey(r.Last)
	if err != nil {
		return "", "", errors.Wrap(err, "getting range end key")
	}
	first, _ := fk.RangeUsing(s.sep)
	_, last := lk.RangeUsing(s.sep)
	return first, last, nil
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
