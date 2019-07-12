package jsonschemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/convert"
)

type testCase struct {
	Name     string
	Schema   string
	Document interface{}
	Assert   struct {
		Match  string
		Errors []validator.ValidationError
	}
}

func TestSchemas(t *testing.T) {

	files, err := filepath.Glob("./tests/*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file // capture
		t.Run(file, func(t *testing.T) {
			a := assert.New(t)

			bs, err := ioutil.ReadFile(file)
			if err != nil {
				a.FailNow("reading case file", err)
			}

			dec := yaml.NewDecoder(bytes.NewReader(bs))

			for {

				tc := testCase{}

				err := dec.Decode(&tc)
				if err == io.EOF {
					break
				}
				if err != nil {
					a.FailNow("unmarshaling case", err)
				}

				t.Run(tc.Name, func(t *testing.T) {
					a := assert.New(t)

					var (
						bs, pbs []byte
						err     error
						vr      *validator.ValidationResult
					)

					bs, err = json.Marshal(cleanupValue(tc.Document))
					if err != nil {
						a.FailNow("marshaling document", err)
					}

					pbs, err = convert.Pretty(bs)
					if err != nil {
						t.Errorf("prettifying document: %v", err)
					}

					t.Logf("validating: %s", string(pbs))
					vr, err = validator.Validate(MustLoadSchema(tc.Schema), bs)
					if err != nil {
						a.FailNow("validating document", err)
					}

					sb := strings.Builder{}
					for _, e := range vr.Errors {
						sb.WriteString("    ")
						sb.WriteString(e.Field)
						sb.WriteString(": ")
						sb.WriteString(e.Description)
						sb.WriteString("\n")
					}
					t.Logf("validation result:\n%s", sb.String())

					missing, lacking := match(tc.Assert.Errors, vr.Errors)
					report := inform(missing, lacking)

					switch tc.Assert.Match {
					case "include":
						if len(missing) > 0 {
							a.Fail("missing expected errors from include match", report)
						}
					case "all":
						if len(missing) > 0 || len(lacking) > 0 {
							a.Fail("expected errors did not match actual", report)
						}
					case "valid":
						if len(lacking) > 0 {
							a.Fail("expected valid document was not", report)
						}
					case "invalid":
						if vr.Valid() {
							a.Fail("expected invalid document was not", report)
						}
					default:
						a.FailNowf("unknown assert match type", "match: %s", tc.Assert.Match)
					}

				})

			}

		})

	}

}

func inform(expected []validator.ValidationError, actual []validator.ValidationError) string {
	s := ""
	if len(expected) > 0 {
		s = fmt.Sprintf("expected errors: %d:\n", len(expected))
		for _, e := range expected {
			s += fmt.Sprintf("%s: %s\n", e.Field, e.Description)
		}
	}
	if len(actual) > 0 {
		s += fmt.Sprintf("actual errors: %d:\n", len(actual))
		for _, a := range actual {
			s += fmt.Sprintf("%s: %s\n", a.Field, a.Description)
		}
	}
	return s
}

func match(expected []validator.ValidationError, actual []validator.ValidationError) (missing []validator.ValidationError, lacking []validator.ValidationError) {
	lacking = actual
	for _, e := range expected {
		found := false
		for i, a := range lacking {
			if e.Field == a.Field && e.Description == a.Description {
				found = true
				lacking = append(lacking[:i], lacking[i+1:]...)
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	return
}

func cleanupValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupArray(v)
	case map[interface{}]interface{}:
		return cleanupMap(v)
	default:
		return v
	}
}

func cleanupArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupValue(v)
	}
	return res
}

func cleanupMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupValue(v)
	}
	return res
}
