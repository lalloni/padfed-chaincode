package schemas

import (
	"io/ioutil"

	"github.com/gobuffalo/packr/v2"
	"github.com/lalloni/gojsonschema"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/convert"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/formats"
)

var fs = packr.New("schemas", "resources")

func init() {
	gojsonschema.Locale = locale{}
	gojsonschema.FormatCheckers.Add("cuit", formats.Cuit)
	gojsonschema.FormatCheckers.Add("periododiario", formats.PeriodoDiario)
	gojsonschema.FormatCheckers.Add("periodomensual", formats.PeriodoMensual)
	gojsonschema.FormatCheckers.Add("periodoanual", formats.PeriodoAnual)
}

func PersonaListSchema() (*gojsonschema.Schema, error) {
	personaSchemaLoader, err := PersonaSchemaJSONLoader()
	if err != nil {
		return nil, errors.Wrap(err, "building loader for persona json schema")
	}
	personaListSchemaloader, err := PersonaListSchemaJSONLoader()
	if err != nil {
		return nil, errors.Wrap(err, "building loader for persona list json schema")
	}
	schemaloader := gojsonschema.NewSchemaLoader()
	schemaloader.Validate = true // validate schema
	schemaloader.Draft = gojsonschema.Draft7
	err = schemaloader.AddSchemas(personaSchemaLoader)
	if err != nil {
		return nil, errors.Wrap(err, "adding referenced json schemas")
	}
	schema, err := schemaloader.Compile(personaListSchemaloader)
	if err != nil {
		return nil, errors.Wrap(err, "building json schema for persona list")
	}
	schema.SetRootSchemaName("(persona list)")
	return schema, nil
}

func PersonaSchema() (*gojsonschema.Schema, error) {
	jsonloader, err := PersonaSchemaJSONLoader()
	if err != nil {
		return nil, errors.Wrap(err, "building loader for persona json schema")
	}
	schemaloader := gojsonschema.NewSchemaLoader()
	schemaloader.Validate = true // validate schema
	schemaloader.Draft = gojsonschema.Draft7
	schema, err := schemaloader.Compile(jsonloader)
	if err != nil {
		return nil, errors.Wrap(err, "building json schema for persona")
	}
	schema.SetRootSchemaName("(persona)")
	return schema, nil
}

func PersonaSchemaJSONLoader() (gojsonschema.JSONLoader, error) {
	return loaderFromYAML("persona.yaml")
}

func PersonaListSchemaJSONLoader() (gojsonschema.JSONLoader, error) {
	return loaderFromYAML("persona-list.yaml")
}

func loaderFromYAML(name string) (gojsonschema.JSONLoader, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, errors.Wrapf(err, "opening '%s'", name)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "reading '%s'", name)
	}
	schema, err := convert.FromYAML(bs, convert.Options{Source: name, Pretty: true})
	if err != nil {
		return nil, errors.Wrapf(err, "converting '%s' to JSON", name)
	}
	loader := gojsonschema.NewBytesLoader(schema)
	_, err = loader.LoadJSON() // for checking json
	if err != nil {
		return nil, errors.Wrapf(err, "parsing JSON converted from '%s'", name)
	}
	return loader, nil
}
