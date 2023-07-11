package main

import (
	"fmt"

	"github.com/demouth/geneva"
)

func main() {

	r := geneva.New()

	v1 := r.Group("/v1", func(c *geneva.Context) {
		fmt.Print("/v1 handler\n")
	})
	v1.GET("/hello", func(c *geneva.Context) {
		c.String(200, "hello v1")
	})

	v2 := r.Group("/v2", func(c *geneva.Context) {
		fmt.Print("/v2 handler\n")
	})
	v2.GET("/hello", func(c *geneva.Context) {
		c.String(200, "hello v2")
	})

	// Listen and Server in 127.0.0.1:8080
	r.Run("127.0.0.1:8080")
}
