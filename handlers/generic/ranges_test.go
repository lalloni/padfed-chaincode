package generic

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestParseRanges(t *testing.T) {
	tests := []struct {
		name   string
		bs     string
		expect interface{}
		err    bool
	}{
		{"point", `bla`,
			queryPoint{key: "bla"},
			false},
		{"singleRange", `[["a","z"]]`,
			[]interface{}{queryRange{begin: "a", until: "z"}},
			false},
		{"doubleRange", `[["a","z"],["1","9"]]`,
			[]interface{}{queryRange{begin: "a", until: "z"}, queryRange{begin: "1", until: "9"}},
			false},
		{"prefixRange", `[["a"]]`,
			[]interface{}{queryRange{begin: "a", until: "a" + string(utf8.MaxRune)}},
			false},
		{"mixedRanges", `[["a"],["b","d"]]`,
			[]interface{}{queryRange{begin: "a", until: "a" + string(utf8.MaxRune)}, queryRange{begin: "b", until: "d"}},
			false},
		{"untilRange", `[["","z"]]`,
			[]interface{}{queryRange{begin: "", until: "z"}},
			false},
		{"beginRange", `[["a",""]]`,
			[]interface{}{queryRange{begin: "a", until: ""}},
			false},
		{"singlePoint", `["a"]`,
			[]interface{}{queryPoint{key: "a"}},
			false},
		{"doublePoint", `["a","b"]`,
			[]interface{}{queryPoint{key: "a"}, queryPoint{key: "b"}},
			false},
		{"mixedPointRange", `["a",["1","9"],"b"]`,
			[]interface{}{queryPoint{key: "a"}, queryRange{begin: "1", until: "9"}, queryPoint{key: "b"}},
			false},
		{"nullIsNotRange", `[null]`,
			nil,
			true},
		{"objectIsNotRange", `[{}]`,
			nil,
			true},
		{"objectIsPoint", `{"a":1}`,
			queryPoint{key: `{"a":1}`},
			false},
		{"stringIsPoint", `"a"`,
			queryPoint{key: `"a"`},
			false},
		{"numberIsPoint", `12`,
			queryPoint{key: "12"},
			false},
		{"emptyPoint", ``,
			queryPoint{key: ""},
			false},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			a := assert.New(t)
			r, err := parseRanges([]byte(test.bs))
			a.EqualValues(test.expect, r)
			if test.err {
				a.Error(err)
			} else {
				a.NoError(err)
			}
		})
	}
}
