package ginger

import (
	"github.com/xeipuuv/gojsonschema"
	"testing"
)

const (
	sampleSchemeName = "sample_scheme"
	sampleScheme     = `{
	"type": "object",
	"additionalProperties": false,
	"required": ["hello", "foo"],
	"properties": {
		"foo": {
			"type": "integer"
		},
		"hello": {
			"type": "string"
		}
	}
}`
)

func init() {
	schemes[sampleSchemeName] = gojsonschema.NewStringLoader(sampleScheme)
}

func TestNewScheme(t *testing.T) {
	Scheme("new-scheme", `{"type":"object"}`)
}

func TestDuplicateScheme(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("duplicate scheme did not panic")
		}
	}()

	Scheme("duplicate", `{"type":"object"}`)
	Scheme("duplicate", `{"type":"object"}`)
}

func TestValidateValid(t *testing.T) {

	if err := validate(map[string]interface{}{
		"hello": "world",
		"foo":   42,
	}, sampleSchemeName); err != nil {
		t.Errorf("fail to validate scheme: %s", err)
	}
}
func TestValidateInvalid(t *testing.T) {

	if err := validate(map[string]interface{}{
		"hello":  "world",
		"foobar": 42,
	}, sampleSchemeName); err == nil {
		t.Errorf("no error on invalid scheme")
	}
}

func TestValidateUnknownScheme(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("validate unknown scheme did not panic")
		}
	}()

	validate(nil, "unknown-scheme")
}

func TestHasScheme(t *testing.T) {

	if !hasScheme(sampleSchemeName) {
		t.Errorf("%s not found (should not)", sampleSchemeName)
	}

	if hasScheme("unknown-scheme") {
		t.Errorf("unkown-scheme found (should not)")
	}
}
