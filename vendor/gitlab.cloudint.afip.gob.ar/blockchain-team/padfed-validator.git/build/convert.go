package build

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/convert"
)

// Convert convierte todos los archivos YAML que se encuentren en source/**/*.yaml
// a archivos JSON ubicados en target/ conservando el path relativo.
func Convert(source, target string) error {
	ss, err := sources(source)
	if err != nil {
		return errors.Wrapf(err, "looking for sources in %s", source)
	}
	for _, s := range ss {
		t := filepath.Join(target, chext(s, ".json"))
		log.Printf("Transforming %s to %s", s, t)
		err := transform(s, t)
		if err != nil {
			return errors.Wrapf(err, "transforming %s", s)
		}
	}
	return nil
}

func chext(s string, ext string) string {
	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s)) + ext
}

func sources(src string) ([]string, error) {
	res := []string{}
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		m, errr := filepath.Match("*.yaml", filepath.Base(path))
		if errr != nil {
			return errr
		}
		if m {
			res = append(res, path)
		}
		return nil
	})
	return res, err
}

func transform(src, tgt string) error {
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	bs, err = convert.FromYAML(bs, convert.Options{
		Source: src,
		Target: tgt,
		Pretty: true,
	})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tgt, bs, 0664)
}
