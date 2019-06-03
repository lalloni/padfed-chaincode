package generic

import (
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/spyzhov/ajson"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
)

func queryRanges(ctx *context.Context, query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case queryPoint:
		s, err := kvget(ctx, q.key)
		if err != nil {
			return nil, errors.Wrap(err, "executing single point query")
		}
		if s.Nil {
			return nil, nil
		}
		return s.Content, nil
	case []interface{}:
		result := []interface{}{}
		for _, q := range q {
			switch q := q.(type) {
			case queryPoint:
				s, err := kvget(ctx, q.key)
				if err != nil {
					return nil, errors.Wrap(err, "executing ranges point query")
				}
				result = append(result, s)
			case queryRange:
				ss, err := krget(ctx, q.begin, q.until)
				if err != nil {
					return nil, errors.Wrap(err, "executing ranges range query")
				}
				result = append(result, ss)
			default:
				return nil, errors.New("internal error")
			}
		}
		return result, nil
	default:
		return nil, errors.New("internal error")
	}
}

func processKeyRanges(ctx *context.Context, query interface{}, process func(key string) (interface{}, error)) (interface{}, error) {
	switch q := query.(type) {
	case queryPoint:
		r, err := process(q.key)
		if err != nil {
			return nil, errors.Wrap(err, "processing single point key")
		}
		return r, nil
	case []interface{}:
		result := []interface{}{}
		for _, q := range q {
			switch q := q.(type) {
			case queryPoint:
				r, err := process(q.key)
				if err != nil {
					return nil, errors.Wrap(err, "processing ranges point key")
				}
				result = append(result, r)
			case queryRange:
				ss, err := krget(ctx, q.begin, q.until)
				if err != nil {
					return nil, errors.Wrap(err, "executing ranges range query")
				}
				rr := []interface{}{}
				for _, s := range ss {
					r, err := process(s.Key)
					if err != nil {
						return nil, errors.Wrap(err, "processing ranges range key")
					}
					rr = append(rr, r)
				}
				result = append(result, rr)
			default:
				return nil, errors.New("internal error")
			}
		}
		return result, nil
	default:
		return nil, errors.New("internal error")
	}
}

func parseRanges(bs []byte) (interface{}, error) {
	root, err := ajson.Unmarshal(bs)
	if err != nil || !root.IsArray() {
		return queryPoint{key: string(bs)}, nil
	}
	query := []interface{}{}
	for _, node := range root.MustArray() {
		switch {
		case node.IsString():
			query = append(query, queryPoint{key: node.MustString()})
		case node.IsArray():
			keys := node.MustArray()
			if !onlyStrings(keys) {
				return nil, errors.Errorf("json query inner arrays must contain strings")
			}
			switch len(keys) {
			case 1:
				a, b := prefixRange(keys[0].MustString())
				query = append(query, queryRange{begin: a, until: b})
			case 2:
				query = append(query, queryRange{begin: keys[0].MustString(), until: keys[1].MustString()})
			default:
				return nil, errors.Errorf("json query inner arrays must contain 1 or 2 elements")
			}
		default:
			return nil, errors.Errorf("json query array must contain strings or string arrays only")
		}
	}
	return query, nil
}

type queryPoint struct {
	key string
}

type queryRange struct {
	begin string
	until string
}

func prefixRange(key string) (string, string) {
	return key, key + string(utf8.MaxRune)
}

func onlyStrings(ns []*ajson.Node) bool {
	for _, n := range ns {
		if !n.IsString() {
			return false
		}
	}
	return true
}
