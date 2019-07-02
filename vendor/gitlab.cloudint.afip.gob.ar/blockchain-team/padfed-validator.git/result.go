package validator

import (
	"strconv"
	"strings"
)

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) String() string {
	l := len(r.Errors)
	switch l {
	case 0:
		return "no errors"
	case 1:
		return r.Errors[0].String()
	default:
		sb := &strings.Builder{}
		sb.WriteString(strconv.Itoa(l))
		sb.WriteString(" errors: ")
		for i, e := range r.Errors {
			sb.WriteString("(" + strconv.Itoa(i+1) + ") ")
			sb.WriteString(e.String())
			if i < l-1 {
				sb.WriteString("; ")
			}
		}
		return sb.String()
	}
}

type ValidationError struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *ValidationError) String() string {
	return e.Field + ": " + e.Description
}
