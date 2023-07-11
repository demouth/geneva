package main

import (
	"log"
	"time"

	"github.com/demouth/geneva"
)

func Logger() geneva.Handler {
	return func(c *geneva.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)
	}
}

func main() {
	r := geneva.New()
	r.Use(Logger())

	r.GET("/test", func(c *geneva.Context) {
		example, _ := c.Get("example")

		// it would print: "12345"
		log.Println(example.(string))
	})

	// Listen and Server in 127.0.0.1:8080
	r.Run("127.0.0.1:8080")
}
