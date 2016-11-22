package ginger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestWrapBadFunc(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Wrap an invalid function did not panic")
		}
	}()

	Wrap(42)
}

func TestWrapBadFuncNoContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Wrap an invalid function did not panic")
		}
	}()

	Wrap(func() {})
}

func TestWrapBadFuncBadInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Wrap an invalid function did not panic")
		}
	}()

	Wrap(func(*gin.Context, map[string]int) {})
}

func TestWrapValidFunction(t *testing.T) {

	Wrap(func(*gin.Context) {})
	Wrap(func(*gin.Context) error { return nil })
	Wrap(func(*gin.Context, map[string]interface{}) {})
	Wrap(func(*gin.Context, map[string]interface{}) error { return nil })
}

func TestErrorHookCall(t *testing.T) {

	err := fmt.Errorf("sample error")

	SetErrorHook(func(c *gin.Context, e error) {
		if e != err {
			t.Errorf("Bad error value")
		}

		err = nil
	})

	fct := Wrap(func(*gin.Context) error { return err })

	fct(&gin.Context{})

	if err != nil {
		t.Errorf("ErrorHook not called")
	}
}

func TestWrapInvalidScheme(t *testing.T) {
	var err error

	SetErrorHook(func(c *gin.Context, e error) {
		err = e
	})

	Scheme("myform", `{
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
	}`)

	fct := Wrap(func(c *gin.Context, data map[string]interface{}) {}, "myform")

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/", bytes.NewBufferString(`{"hello":"world"}`))
	req.Header.Set("Content-Type", "application/json")

	c := &gin.Context{
		Request: req,
	}

	fct(c)

	if err == nil {
		t.Errorf("ErrorHook not called")
	}
}

func TestWrapValidScheme(t *testing.T) {
	var err error

	SetErrorHook(func(c *gin.Context, e error) {
		err = e
	})

	Scheme("myform2", `{
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
	}`)

	fct := Wrap(func(c *gin.Context, data map[string]interface{}) {}, "myform2")

	req, _ := http.NewRequest(http.MethodPost, "http://localhost/", bytes.NewBufferString(`{"hello":"world", "foo": 42}`))
	req.Header.Set("Content-Type", "application/json")

	c := &gin.Context{
		Request: req,
	}

	fct(c)

	if err != nil {
		t.Errorf("ErrorHook called")
	}
}
