# Geneva Web Framework

Geneva is a simple Gin-like web framework written in Golang.

## Running Geneva

```golang
package main

import (
	"net/http"

	"github.com/demouth/geneva"
)

func main() {
	r := geneva.New()
	v1 := r.Group("/v1")
	v1.GET("/ping", func(c *geneva.Context) {
		c.JSON(http.StatusOK, geneva.H{
			"message": "pong",
		})
	})
	r.Run("127.0.0.1:8080")
}
```