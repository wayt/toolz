package ginger

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/xeipuuv/gojsonschema"
)

// schemes store
var schemes = make(map[string]gojsonschema.JSONLoader)

// Scheme register a new jsonschema
// It uses https://github.com/xeipuuv/gojsonschema
// Check for jsonschema doc http://json-schema.org/
func Scheme(name, jsonScheme string) gojsonschema.JSONLoader {
	if _, ok := schemes[name]; ok {
		panic(fmt.Errorf("duplicate scheme name: %s", name))
	}

	loader := gojsonschema.NewStringLoader(jsonScheme)
	schemes[name] = loader

	return loader
}

// validate a data input for a givent scheme
func validate(data map[string]interface{}, schemeName string) error {

	schema, ok := schemes[schemeName]
	if !ok {
		panic(fmt.Errorf("invalid scheme name: %s", schemeName))
	}

	doc := gojsonschema.NewGoLoader(data)

	if res, err := gojsonschema.Validate(schema, doc); err != nil {
		return errors.Annotate(err, "fail to validate jsonschema")
	} else if !res.Valid() {
		return errors.BadRequestf("bad %s value", res.Errors()[0].Field())
	}

	return nil
}

// hasScheme verify a scheme exist
func hasScheme(name string) bool {
	_, ok := schemes[name]
	return ok
}
