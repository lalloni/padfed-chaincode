package validator

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/lalloni/gojsonschema"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/convert"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/formats"
)

var fs = packr.New("schemas", "./schemas")

func Schemas() []string {
	names := []string(nil)
	for _, f := range fs.List() {
		if !strings.HasSuffix(f, ".yaml") {
			continue
		}
		names = append(names, strings.TrimSuffix(f, filepath.Ext(f)))
	}
	return names
}

func MustLoadSchema(name string) *gojsonschema.Schema {
	s, err := LoadSchema(name)
	if err != nil {
		panic(err)
	}
	return s
}

func LoadSchema(name string) (*gojsonschema.Schema, error) {

	var root gojsonschema.JSONLoader
	loaders := []gojsonschema.JSONLoader(nil)

	for _, f := range fs.List() {
		if !strings.HasSuffix(f, ".yaml") {
			continue
		}
		loader, err := loaderFromYAML(f)
		if err != nil {
			return nil, errors.Wrapf(err, "creating json loader for %q", f)
		}
		if f == name+".yaml" {
			root = loader
		} else {
			loaders = append(loaders, loader)
		}
	}

	if root == nil {
		return nil, errors.Errorf("schema not found: %s", name)
	}

	schemaloader := gojsonschema.NewSchemaLoader()
	schemaloader.Validate = true // validate schema
	schemaloader.Draft = gojsonschema.Draft7

	err := schemaloader.AddSchemas(loaders...)
	if err != nil {
		return nil, errors.Wrap(err, "adding schemas")
	}

	schema, err := schemaloader.Compile(root)
	if err != nil {
		return nil, errors.Wrapf(err, "building json schema for %q", name)
	}

	schema.SetRootSchemaName("(" + name + ")")

	return schema, nil

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

func init() {
	gojsonschema.Locale = locale{}
	gojsonschema.FormatCheckers.Add("cuit", formats.Cuit)
	gojsonschema.FormatCheckers.Add("periododiario", formats.PeriodoDiario)
	gojsonschema.FormatCheckers.Add("periodomensual", formats.PeriodoMensual)
	gojsonschema.FormatCheckers.Add("periodoanual", formats.PeriodoAnual)
}
