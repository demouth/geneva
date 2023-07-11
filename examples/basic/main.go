package main

import (
	"fmt"

	"github.com/demouth/geneva"
)

func main() {
	r := geneva.New()

	// Parameters in path
	r.GET("/hello/:name", hello)

	// Using HTTP method
	r.GET("/greeting", greeting)
	r.POST("/greeting", greeting)
	r.PUT("/greeting", greeting)
	r.DELETE("/greeting", greeting)
	r.PATCH("/greeting", greeting)
	r.HEAD("/greeting", greeting)
	r.OPTIONS("/greeting", greeting)

	// Query string parameters
	r.GET("/welcome", welcome)

	// Query string and post form
	r.POST("/post", post)

	// Listen and Server in 127.0.0.1:8080
	r.Run("127.0.0.1:8080")
}

func hello(c *geneva.Context) {
	s := fmt.Sprintf("hello, %s!", c.Param("name"))
	c.String(200, s)
}

func greeting(c *geneva.Context) {
	c.String(200, c.Request.Method)
}

// Query string parameters are parsed using the existing underlying request object.
// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
func welcome(c *geneva.Context) {
	firstname := c.DefaultQuery("firstname", "Guest")
	lastname := c.Query("lastname")
	s := fmt.Sprintf("Hello %s %s", firstname, lastname)
	c.String(200, s)
}

// Query string and post form
//
// cURL:
//
//	curl "127.0.0.1:8080/post?id=1234" \
//	-d "name=manu" \
//	-d "message=this_is_great"
//
// HTTP:
//
//	POST /post?id=1234 HTTP/1.1
//	Host: 127.0.0.1:8080
//	Content-Type: application/x-www-form-urlencoded
//
//	name=manu&message=this_is_great
//
// RESPONSE:
//
//	id: 1234; page: 0; name: manu; message: this_is_great
func post(c *geneva.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	s := fmt.Sprintf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)

	c.String(200, s)
}
