package convert

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/bitly/go-simplejson"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
)

type Options struct {
	Source string
	Target string
	Pretty bool
}

func FromYAML(bs []byte, opts Options) ([]byte, error) {
	bs, err := yaml.YAMLToJSON(bs)
	if err != nil {
		return nil, err
	}
	bs, err = clean(bs)
	if err != nil {
		return nil, err
	}
	bs, err = patch(bs, opts)
	if err != nil {
		return nil, err
	}
	if opts.Pretty {
		bs, err = Pretty(bs)
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}

const p = `
{
	"$meta": {
		{{ if .Source }}"source": "{{ .Source }}",{{ end }}
		"comment": "SCHEMA GENERADO AUTOMÁTICAMENTE (NO MODIFICAR)"
	}
}`

var tpl = template.Must(template.New("patch").Funcs(sprig.TxtFuncMap()).Parse(p))

func patch(bs []byte, opts Options) ([]byte, error) {
	b := bytes.Buffer{}
	err := tpl.Execute(&b, opts)
	if err != nil {
		return nil, err
	}
	return jsonpatch.MergePatch(bs, b.Bytes())
}

func clean(bs []byte) ([]byte, error) {
	v, err := simplejson.NewJson(bs)
	if err != nil {
		return nil, err
	}
	filter(v)
	return v.Encode()
}

func filter(v *simplejson.Json) {
	m, err := v.Map()
	if err != nil {
		return
	}
	for k := range m {
		if strings.HasPrefix(k, "x-") {
			v.Del(k)
		} else {
			filter(v.Get(k))
		}
	}
}

func Pretty(bs []byte) ([]byte, error) {
	b := bytes.Buffer{}
	err := json.Indent(&b, bs, "", "  ")
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
