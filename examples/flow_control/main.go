package main

import (
	"github.com/demouth/geneva"
)

func main() {
	r := geneva.New()

	r.Handle(
		"GET",
		"/hello",
		func(c *geneva.Context) {
			// do something
		},
		func(c *geneva.Context) {
			// something error happens
			if true {
				c.AbortWithStatus(500)
			}
		},
		func(c *geneva.Context) {
			// never called
			c.String(200, "hello")
		},
	)

	// Listen and Server in 127.0.0.1:8080
	r.Run("127.0.0.1:8080")
}
