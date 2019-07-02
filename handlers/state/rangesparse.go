package state

import (
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/spyzhov/ajson"
)

func Parse(bs []byte) (*Ranges, error) {
	root, err := ajson.Unmarshal(bs)
	if err != nil || !root.IsArray() {
		return Single(Point(string(bs))), nil
	}
	query := []*Item{}
	for _, node := range root.MustArray() {
		switch {
		case node.IsString():
			query = append(query, Point(node.MustString()))
		case node.IsArray():
			keys := node.MustArray()
			if !onlyStrings(keys) {
				return nil, errors.Errorf("json query inner arrays must contain strings")
			}
			switch len(keys) {
			case 1:
				a, b := prefixRange(keys[0].MustString())
				query = append(query, Range(a, b))
			case 2:
				query = append(query, Range(keys[0].MustString(), keys[1].MustString()))
			default:
				return nil, errors.Errorf("json query inner arrays must contain 1 or 2 elements")
			}
		default:
			return nil, errors.Errorf("json query array must contain strings or string arrays only")
		}
	}
	return List(query...), nil
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
