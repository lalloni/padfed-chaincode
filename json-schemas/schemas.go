package jsonschemas

import (
	packr "github.com/gobuffalo/packr/v2"
	"github.com/lalloni/gojsonschema"

	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

var fs = packr.New("schemas", "./schemas/")

func Schemas() []string {
	return validator.Schemas(fs)
}

func MustLoadSchema(name string) *gojsonschema.Schema {
	return validator.MustLoadSchema(fs, name)
}

func LoadSchema(name string) (*gojsonschema.Schema, error) {
	return validator.LoadSchema(fs, name)
}
