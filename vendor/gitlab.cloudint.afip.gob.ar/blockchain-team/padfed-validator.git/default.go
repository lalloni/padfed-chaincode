package validator

import "fmt"

var Default = MustNew()

func MustNew() Validator {
	v, err := New()
	if err != nil {
		panic(fmt.Sprintf("Error creating default validator: %v", err))
	}
	return v
}

func ValidatePersonaJSON(bs []byte) (*ValidationResult, error) {
	return Default.ValidatePersonaJSON(bs)
}

func ValidatePersonaListJSON(bs []byte) (*ValidationResult, error) {
	return Default.ValidatePersonaListJSON(bs)
}
