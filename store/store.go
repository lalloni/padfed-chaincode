package store

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/filtering"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/key"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/marshaling"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/store/meta"
)

type Store interface {
	PutValue(key *key.Key, val interface{}) error
	GetValue(key *key.Key, val interface{}) (bool, error)
	HasValue(key *key.Key) (bool, error)
	DelValue(key *key.Key) error

	PutComposite(com *meta.PreparedComposite, val interface{}) error
	GetComposite(com *meta.PreparedComposite, id interface{}) (interface{}, error)
	HasComposite(com *meta.PreparedComposite, id interface{}) (bool, error)
	DelComposite(com *meta.PreparedComposite, id interface{}) error
}

func New(stub shim.ChaincodeStubInterface, opts ...Option) Store {
	s := &simplestore{
		stub:       stub,
		marshaling: DefaultMarshaling,
		filtering:  DefaultFiltering,
		sep:        key.DefaultSep,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type simplestore struct {
	stub       shim.ChaincodeStubInterface
	marshaling marshaling.Marshaling
	filtering  filtering.Filtering
	sep        *key.Sep
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
	if err := s.internalPutValue(we.Key, we.Value); err != nil {
		return errors.Wrapf(err, "putting composite %q witness", com.Name)
	}
	for _, entry := range com.SingletonsEntries(val) {
		if entry.Value != nil {
			if err := s.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q singleton %q", com.Name, entry)
			}
		}
	}
	for _, entry := range com.CollectionsEntries(val) {
		if entry.Value == nil {
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
	val := com.CreateValue()
	com.SetIdentifier(val, id)
	start, end := com.RangeSep(valkey, s.sep)
	states, err := s.stub.GetStateByRange(start, end)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q with key %q states iterator", com.Name, valkey)
	}
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
				return nil, errors.Wrapf(err, "parsing composite %q with key %q collection item %q value", com.Name, valkey, statekey)
			}
			member.Collector(val, meta.Item{Identifier: statekey.Tag.Value, Value: itemval})
		case com.Singleton(statekey) != nil:
			member := com.Singleton(statekey)
			itemval := member.Creator()
			err := s.internalParseValue(state.GetValue(), itemval)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing composite %q with key %q singleton item %q value", com.Name, valkey, statekey)
			}
			member.Setter(val, itemval)
		}
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
	start, end := com.RangeSep(key, s.sep)
	states, err := s.stub.GetStateByRange(start, end)
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
