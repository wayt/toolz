# Toolz

This is a set of Golang tools designed to ease webapp creation

## Ginger

Ginger is a gin (https://github.com/gin-gonic/gin) wrapper, that allow you to return errors, define a global error handler and validate input form using json schema

### Start using it

1. Install it

    ```sh
    $ go get github.com/wayt/toolz/ginger
    ```

2. Import it

    ```go
    import "github.com/wayt/toolz/ginger"
    ```

3. (Optional) Import `net/http`. This is required for example if using constants such as `http.StatusOK`.

    ```go
    import "net/http"
    ```

### Examples

#### Basic GET with error

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wayt/toolz/ginger"
	"net/http"
)

func main() {

	r := gin.Default()

	ginger.SetErrorHook(func(c *gin.Context, e error) {
		c.String(http.StatusInternalServerError, e.Error())
	})

	r.GET("/hello", ginger.Wrap(func(c *gin.Context) (err error) {
        err = fmt.Errorf("Hello World !")
        return
    }))

	r.Run()
}
```

#### Post with json schema

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wayt/toolz/ginger"
	"net/http"
)

var _ = ginger.Scheme("myform", `{
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

func main() {

	r := gin.Default()

	ginger.SetErrorHook(func(c *gin.Context, e error) {
		c.String(http.StatusInternalServerError, e.Error())
	})

	r.POST("/", ginger.Wrap(func(c *gin.Context) (err error) {
        c.JSON(http.StatusOK, data)
        return
    }, "myform"))

	r.Run()
}
```
