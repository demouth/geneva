package geneva

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestString(t *testing.T) {
	// SETUP
	r := New()
	r.Handle("GET", "/test", func(c *Context) {
		c.String(200, "test")
	})

	// RUN
	w := request(r, "GET", "/test")

	// TEST
	if w.Code != 200 {
		t.Errorf("Status code should be %v, was %d", 200, w.Code)
	}
	if w.Body.String() != "test" {
		t.Errorf("Body should be %s, was %s", "test", w.Body.String())
	}
	if w.HeaderMap["Content-Type"][0] != "text/plain" {
		t.Errorf("Content-Type should be %s, was %s", "text/plain", w.HeaderMap["Content-Type"])
	}
}

func TestJSON(t *testing.T) {
	// SETUP
	r := New()
	r.Handle("GET", "/test", func(c *Context) {
		c.JSON(200, H{"foo": "bar"})
	})

	// RUN
	w := request(r, "GET", "/test")

	// TEST
	if w.Code != 200 {
		t.Errorf("Status code should be %v, was %d", 200, w.Code)
	}
	if w.Body.String() != "{\"foo\":\"bar\"}\n" {
		t.Errorf("Body should be {\"foo\":\"bar\"}, was %s", w.Body.String())
	}
	if w.HeaderMap["Content-Type"][0] != "application/json" {
		t.Errorf("Content-Type should be %s, was %s", "application/json", w.HeaderMap["Content-Type"])
	}
}

func TestParams(t *testing.T) {
	// SETUP
	r := New()
	r.Handle("POST", "/test/:name", func(c *Context) {
		if name := c.Param("name"); name != "geneva" {
			t.Errorf("name = %q; want \"geneva\"", name)
		}
		if v := c.Query("q1"); v != "1" {
			t.Errorf("q1 = %q; want 1", v)
		}
		if v := c.Query("q2"); v != "2" {
			t.Errorf("q2 = %q; want 2", v)
		}
		if v := c.Query("q3"); v != "" {
			t.Errorf("q3 = %q; want \"\"", v)
		}
		if v := c.DefaultPostForm("p3", "default"); v != "default" {
			t.Errorf("p3  = %q; want \"default\"", v)
		}
		if v := c.PostForm("p1"); v != "1" {
			t.Errorf("p1  = %q; want \"1\"", v)
		}
		if v := c.PostForm("p2"); v != "2" {
			t.Errorf("p2  = %q; want \"2\"", v)
		}
	})

	// RUN
	w := postFormRequest(r, "/test/geneva?q1=1&q2=2", "p1=1&p2=2&p1=3")

	if w.Code != 200 {
		t.Errorf("Response code should be Bad request, was: %d", w.Code)
	}
}

func TestAbort(t *testing.T) {

	i := 0

	r := New()
	r.Use(Recovery())
	r.GET("/abort",
		func(c *Context) {
			i++
		},
		func(c *Context) {
			i++
			c.AbortWithStatus(501)
			i++
		},
		func(c *Context) {
			i++
		},
	)

	w := request(r, "GET", "/abort")

	if i != 3 {
		t.Errorf("`i` should be 3, was %d", i)
	}

	if w.Code != 501 {
		t.Errorf("Response code should be Bad request, was: %d", w.Code)
	}
}

func TestRecovery(t *testing.T) {

	r := New()
	r.Use(Recovery())
	r.GET("/recovery", func(c *Context) {
		panic("something problem")
	})

	w := request(r, "GET", "/recovery")

	if w.Code != 500 {
		t.Errorf("Response code should be Bad request, was: %d", w.Code)
	}
}

func TestContext(t *testing.T) {

	i := 0

	// SETUP
	r := New()
	r.Use(Recovery())
	r.Use(
		func(c *Context) {
			_, exists := c.Get("unknown")
			if exists {
				t.Errorf("Exists should be %t, was %t", false, exists)
			}
		},
		func(c *Context) {
			c.Set("0", i)
			i++
		},
		func(c *Context) {
			c.Set("1", i)
			i++
		},
	)
	r.Use(
		func(c *Context) {
			c.Set("2", i)
			i++
		},
	)
	group := r.Group(
		"/group",
		func(c *Context) {
			c.Set("3", i)
			i++
		},
		func(c *Context) {
			c.Set("4", i)
			i++
		},
	)
	group.Use(
		func(c *Context) {
			c.Set("5", i)
			i++
		},
	)
	group.Handle(
		"GET",
		"/handle",
		func(c *Context) {
			c.Set("6", i)
			i++
		},
		func(c *Context) {
			c.Set("7", i)
			i++
		},
		func(c *Context) {
			if v, _ := c.Get("0"); v != 0 {
				t.Errorf("Get value should be %d, was %d", 0, v)
			}
			if v, _ := c.Get("1"); v != 1 {
				t.Errorf("Get value should be %d, was %d", 1, v)
			}
			if v, _ := c.Get("2"); v != 2 {
				t.Errorf("Get value should be %d, was %d", 2, v)
			}
			if v, _ := c.Get("3"); v != 3 {
				t.Errorf("Get value should be %d, was %d", 3, v)
			}
			if v, _ := c.Get("4"); v != 4 {
				t.Errorf("Get value should be %d, was %d", 4, v)
			}
			if v, _ := c.Get("5"); v != 5 {
				t.Errorf("Get value should be %d, was %d", 5, v)
			}
			if v, _ := c.Get("6"); v != 6 {
				t.Errorf("Get value should be %d, was %d", 6, v)
			}
			if v, _ := c.Get("7"); v != 7 {
				t.Errorf("Get value should be %d, was %d", 7, v)
			}
		},
	)

	// RUN
	w := request(r, "GET", "/group/handle")

	// TEST
	if w.Code != 200 {
		t.Errorf("Status code should be %v, was %d", 200, w.Code)
	}
}

func request(e *Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w
}

// "foo=bar&page=11&both=&foo=second"
func postFormRequest(e *Engine, path, body string) *httptest.ResponseRecorder {
	bf := bytes.NewBufferString(body)
	req, _ := http.NewRequest("POST", path, bf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w
}
