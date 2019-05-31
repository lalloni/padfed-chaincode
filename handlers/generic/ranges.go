package generic

import (
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/spyzhov/ajson"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func processRanges(ctx *context.Context, ranges []byte, processor func(single bool, key string) (interface{}, error)) *response.Response {
	q, err := parseRanges(ranges)
	if err != nil {
		return response.BadRequest("invalid argument: %v", err)
	}
	switch q := q.(type) {
	case queryPoint:
		v, err := processor(true, q.key)
		if err != nil {
			return response.Error("processing range key: %v", err)
		}
		if res, ok := v.(*response.Response); ok {
			return res
		}
		return response.OK(v)
	case []interface{}:
		result := []interface{}{}
		for _, q := range q {
			switch q := q.(type) {
			case queryPoint:
				v, err := processor(false, q.key)
				if err != nil {
					return response.Error("processing range key: %v", err)
				}
				if res, ok := v.(*response.Response); ok {
					return res
				}
				result = append(result, v)
			case queryRange:
				ks, res := rangekeys(ctx, q.begin, q.until)
				if res != nil {
					return res
				}
				rr := []interface{}{}
				for _, k := range ks {
					v, err := processor(false, k)
					if err != nil {
						return response.Error("processing range key: %v", err)
					}
					if res, ok := v.(*response.Response); ok {
						return res
					}
					rr = append(rr, v)
				}
				result = append(result, rr)
			default:
				return response.Error("internal error")
			}
		}
		return response.OK(result)
	default:
		return response.Error("internal error")
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
