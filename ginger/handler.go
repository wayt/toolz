package ginger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

var (
	contextType            = reflect.TypeOf((*gin.Context)(nil))
	errorType              = reflect.TypeOf((*error)(nil)).Elem()
	mapStringInterfaceType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
	errHookFct             = defaultErrHookFct
)

func Wrap(f interface{}, schema ...string) gin.HandlerFunc {

	fv := reflect.ValueOf(f)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		panic("not a function")
	}

	if ft.NumIn() == 0 || ft.In(0) != contextType {
		panic("first func argument must be a *gin.Context")
	}

	hasIn := ft.NumIn() > 1
	if hasIn {
		if ft.In(1) != mapStringInterfaceType {
			panic("second func argument must be a map[string]interface{}")
		}

		if len(schema) > 0 {
			if !hasScheme(schema[0]) {
				panic(fmt.Errorf("unknown schema: %s", schema[0]))
			}
		}
	}

	return func(c *gin.Context) {

		in := []reflect.Value{reflect.ValueOf(c)}

		if hasIn {

			data := make(map[string]interface{})

			if err := c.Bind(&data); err != nil {
				errHookFct(c, err)
				return
			}

			if len(schema) > 0 {
				if err := validate(data, schema[0]); err != nil {
					errHookFct(c, err)
					return
				}
			}

			in = append(in, reflect.ValueOf(data))
		}

		out := fv.Call(in)

		if ft.NumOut() == 1 && ft.Out(0) == errorType {
			if errv := out[0]; !errv.IsNil() {
				errHookFct(c, errv.Interface().(error))
				return
			}
		}
	}
}

type ErrorHook func(*gin.Context, error)

func SetErrorHook(e ErrorHook) {
	errHookFct = e
}

func defaultErrHookFct(c *gin.Context, e error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": e.Error(),
	})
}
