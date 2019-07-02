package state

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		bs     string
		expect interface{}
		err    bool
	}{
		{"point", `bla`,
			Single(Point("bla")),
			false},
		{"singleRange", `[["a","z"]]`,
			List(Range("a", "z")),
			false},
		{"doubleRange", `[["a","z"],["1","9"]]`,
			List(Range("a", "z"), Range("1", "9")),
			false},
		{"prefixRange", `[["a"]]`,
			List(Range("a", "a"+string(utf8.MaxRune))),
			false},
		{"mixedRanges", `[["a"],["b","d"]]`,
			List(Range("a", "a"+string(utf8.MaxRune)), Range("b", "d")),
			false},
		{"untilRange", `[["","z"]]`,
			List(Range("", "z")),
			false},
		{"beginRange", `[["a",""]]`,
			List(Range("a", "")),
			false},
		{"singlePoint", `["a"]`,
			List(Point("a")),
			false},
		{"doublePoint", `["a","b"]`,
			List(Point("a"), Point("b")),
			false},
		{"mixedPointRange", `["a",["1","9"],"b"]`,
			List(Point("a"), Range("1", "9"), Point("b")),
			false},
		{"nullIsNotRange", `[null]`,
			nil,
			true},
		{"objectIsNotRange", `[{}]`,
			nil,
			true},
		{"objectIsPoint", `{"a":1}`,
			Single(Point(`{"a":1}`)),
			false},
		{"stringIsPoint", `"a"`,
			Single(Point(`"a"`)),
			false},
		{"numberIsPoint", `12`,
			Single(Point("12")),
			false},
		{"emptyPoint", ``,
			Single(Point("")),
			false},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			a := assert.New(t)
			r, err := Parse([]byte(test.bs))
			if test.err {
				a.Error(err)
			} else {
				a.EqualValues(test.expect, r)
				a.NoError(err)
			}
		})
	}
}
