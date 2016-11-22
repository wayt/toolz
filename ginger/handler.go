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

// Wrap overload gin.HandlerFunc to add few features
// Handler func first argument MUST be a *gin.Context
// Your handler func can return an error, for example:
//		gin.GET("/", ginger.Wrap(func(c *gin.Context) error {
//          if c.Query("email") == "" {
//              return fmt.Errorf("missing email")
//          }
//          return nil
//      })
// It will be passed to the ErrorHook handler, See SetErrorHook
// It also support form validation using json schema, first create
// a ginger.Scheme, then specify its name when you register your handler
// for example:
//		gin.POST("/", ginger.Wrap(func(c *gin.Context, data map[string]interface{}) {
//          // Your form is in data
//          // To process it, use https://github.com/mitchellh/mapstructure
//      })
// Your handler func MUST have a map[string]interface{} as 2nd argument
func Wrap(f interface{}, schema ...string) gin.HandlerFunc {

	fv := reflect.ValueOf(f)
	ft := fv.Type()

	// Ensure f is a function
	if ft.Kind() != reflect.Func {
		panic("not a function")
	}

	// Ensure first argument is a *gin.Context
	if ft.NumIn() == 0 || ft.In(0) != contextType {
		panic("first func argument must be a *gin.Context")
	}

	hasIn := ft.NumIn() > 1
	if hasIn {
		// If we have a 2nd argument, ensure it's a map[string]interface{}
		if ft.In(1) != mapStringInterfaceType {
			panic("second func argument must be a map[string]interface{}")
		}

		if len(schema) > 0 {
			// Ensure scheme is already registered, See ginger.Scheme()
			if !hasScheme(schema[0]) {
				panic(fmt.Errorf("unknown schema: %s", schema[0]))
			}
		}
	}

	// Build the gin.HandlerFunc
	return func(c *gin.Context) {

		// First arg will always be the gin Context
		in := []reflect.Value{reflect.ValueOf(c)}

		if hasIn {
			// Process input form

			data := make(map[string]interface{})

			if err := c.Bind(&data); err != nil {
				errHookFct(c, err)
				return
			}

			// Validate form if needed
			if len(schema) > 0 {
				if err := validate(data, schema[0]); err != nil {
					errHookFct(c, err)
					return
				}
			}

			in = append(in, reflect.ValueOf(data))
		}

		// The actual call
		out := fv.Call(in)

		// Check for output error
		if ft.NumOut() == 1 && ft.Out(0) == errorType {
			if errv := out[0]; !errv.IsNil() {
				errHookFct(c, errv.Interface().(error))
				return
			}
		}
	}
}

// ErrorHook defaine error function handler
type ErrorHook func(*gin.Context, error)

// SetErrorHook change the appliction error handler
// it will be called when Wrap encouter an error or
// an error is return from the HandlerFunc
func SetErrorHook(e ErrorHook) {
	errHookFct = e
}

// defaultErrHookFct is the default error handler
func defaultErrHookFct(c *gin.Context, e error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": e.Error(),
	})
}
