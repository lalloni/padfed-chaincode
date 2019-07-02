// MIT License
//
// Copyright (c) 2019 Pablo Ignacio Lalloni
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cuit

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

const (
	// Min is the minimum valid cuit
	Min = 20000000001

	// Max is the maximum valid cuit
	Max = 34999999990
)

var (
	legkinds = []uint64{20, 24, 27, 30, 34}
	allkinds = []uint64{20, 23, 24, 27, 30, 33, 34}
	factor   = []uint64{2, 3, 4, 5, 6, 7}
	factors  = len(factor)
	pattern  = regexp.MustCompile(`^(\d{2})-?(\d{8})-?(\d{1})$`)
)

// IsValid checks the provided CUIT/CUIL number for validity considering
// wether the number is in valid range, has a valid kind (20, 23, 24, 27,
// 30, 33 or 34) and has a correct verifier digit.
func IsValid(cuit uint64) bool {
	return validSize(cuit) && validKind(cuit) && validVerifier(cuit)
}

func validSize(cuit uint64) bool {
	return cuit/1e11 == 0
}

func validKind(cuit uint64) bool {
	cuitk := (cuit % 1e11) / 1e9
	valid := false
	for _, kind := range allkinds {
		if cuitk == kind {
			valid = true
			break
		}
	}
	return valid
}

func validVerifier(cuit uint64) bool {
	return cuit%10 == Verifier(cuit)
}

// Verifier computes and returns the correct verifier digit for the
// number supplied in the cuit parameter.
//
// This function ignores the unit digit from the input value (as it is the
// verifier digit being calculated).
//
// This function ignores digits past the eleventh (10^11).
//
// A return value of 10 means the input cuit number does not exist.
func Verifier(cuit uint64) uint64 {
	var num uint64
	rem := (cuit % 1e11) / 10 // drop out of range and verifier digits
	for i := 0; rem != 0; i++ {
		num += factor[i%factors] * (rem % 10)
		rem /= 10
	}
	num = 11 - num%11
	if num == 11 {
		return 0
	}
	return num
}

// Parse extracts a CUIT number from the string provided using the
// standard format "DD-DDDDDDDD-D" being "D" any decimal digit and
// both "-" characters optional.
//
// If the string can not be parsed as a CUIT number the function returns error.
func Parse(cuit string) (uint64, error) {
	match := pattern.FindStringSubmatch(cuit)
	if match == nil {
		return 0, errors.Errorf("formato incorrecto de cuit/cuil: %q", cuit)
	}
	// the following three errors can never
	// happen because the regexp pattern
	// ensures all submatches are **only**
	// digits, so we can safely ignore them
	kind, _ := strconv.Atoi(match[1])
	id, _ := strconv.Atoi(match[2])
	ver, _ := strconv.Atoi(match[3])
	return uint64(kind)*1e9 + uint64(id)*1e1 + uint64(ver), nil
}

// Parts extracts the kind number, identifier number and verifier digit
// parts out of the provided CUIT number and returns them as the three
// return values.
//
// This functions discards any decimal digit exceeding the allowed range.
func Parts(cuit uint64) (kind, id, ver uint64) {
	clean := cuit % 1e11 // get rid of possible excess digits
	return clean / 1e9, (clean % 1e9) / 10, clean % 10
}

// Format returns a standard formatted string representation of the provided
// CUIT number.
//
// This function uses Parts(cuit) to extract the constituent parts of the
// CUIT number, hence that function behavior regarding digits exceeding
// the allowed range is maintained.
func Format(cuit uint64) string {
	kind, id, ver := Parts(cuit)
	return fmt.Sprintf("%02d-%08d-%01d", kind, id, ver)
}

// Random computes and returns random valid cuit numbers.
func Random(r *rand.Rand) uint64 {
	k := legkinds[r.Intn(len(legkinds))]
	id := r.Uint64() % 1e8
	c := Compose(k, id, 0)
	v := Verifier(c)
	if v == 10 {
		if k < 30 {
			c = Compose(23, id, 0)
		} else {
			c = Compose(33, id, 0)
		}
		v = Verifier(c)
	}
	return c + v
}

// Compose builds a cuit number from its parts.
//
// This function drops the most-significative excess digits of all its input
// arguments.
//
// Please see the examples.
func Compose(kind, id, ver uint64) uint64 {
	return (kind%1e2)*1e9 + (id%1e8)*10 + ver%1e1
}

// Pred returns predecessor of cuit unless cuit equals Min in which
// case it returns Min or cuit is not valid in which case it returns 0.
func Pred(cuit uint64) uint64 {
	if !IsValid(cuit) {
		return 0
	}
	if cuit == Min {
		return cuit
	}
	pre, num, _ := Parts(cuit)
	r := pre*1e9 + num*1e1
	for {
		r -= 10
		v := Verifier(r)
		if v < 10 {
			return r + v
		}
	}
}

// Succ returns the successor of cuit unless cuit equals Max in which
// case it returns Max or cuit is not valid in which case it returns 0.
func Succ(cuit uint64) uint64 {
	if !IsValid(cuit) {
		return 0
	}
	if cuit == Max {
		return cuit
	}
	pre, num, _ := Parts(cuit)
	r := pre*1e9 + num*1e1
	for {
		r += 10
		v := Verifier(r)
		if v < 10 {
			return r + v
		}
	}
}
