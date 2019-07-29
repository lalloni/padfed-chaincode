package jsonschemas

import (
	"github.com/lalloni/gojsonschema"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/resources"

	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
)

func Schemas() []string {
	return validator.Schemas(resources.Schemas)
}

func MustLoadSchema(name string) *gojsonschema.Schema {
	return validator.MustLoadSchema(resources.Schemas, name)
}

func LoadSchema(name string) (*gojsonschema.Schema, error) {
	return validator.LoadSchema(resources.Schemas, name)
}
