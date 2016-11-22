package ginger

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/xeipuuv/gojsonschema"
)

var forms = make(map[string]gojsonschema.JSONLoader)

func Scheme(name, jsonScheme string) gojsonschema.JSONLoader {
	if _, ok := forms[name]; ok {
		panic(fmt.Errorf("duplicate scheme name: %s", name))
	}

	loader := gojsonschema.NewStringLoader(jsonScheme)
	forms[name] = loader

	return loader
}

func validate(data map[string]interface{}, formName string) error {

	schema, ok := forms[formName]
	if !ok {
		panic(fmt.Errorf("invalid form name: %s", formName))
	}

	doc := gojsonschema.NewGoLoader(data)

	if res, err := gojsonschema.Validate(schema, doc); err != nil {
		return errors.Annotate(err, "fail to validate jsonschema")
	} else if !res.Valid() {
		return errors.BadRequestf("bad %s value", res.Errors()[0].Field())
	}

	return nil
}

func hasScheme(name string) bool {
	_, ok := forms[name]
	return ok
}
