// Package geneva is a simple web framework.
package geneva

import (
	"encoding/json"
	"math"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

const (
	abortIndex = math.MaxInt >> 1
)

type (
	// Engine is the framework's instance.
	Engine struct {
		*RouterGroup
		router *httprouter.Router
	}
	// RouterGroup is used internally to configure router.
	RouterGroup struct {
		path     string
		handlers Handlers
		engine   *Engine
	}
	// Context is argument of Handler.
	Context struct {
		Writer   http.ResponseWriter
		Request  *http.Request
		Params   httprouter.Params
		keys     map[string]any
		handlers Handlers
		index    int
	}
	// Handler defines the handler used by geneva middleware as return value.
	Handler func(*Context)
	// Handlers defines a Handler slice.
	Handlers []Handler
	// H is a shortcut for map[string]any
	H map[string]any
)

////////////
// Engine //
////////////

// New returns a new Root instance.
func New() *Engine {
	e := &Engine{
		RouterGroup: &RouterGroup{
			path: "",
		},
		router: httprouter.New(),
	}
	e.RouterGroup.engine = e
	return e
}

// Run listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}

/////////////////
// RouterGroup //
/////////////////

// GET is a shortcut for geneva.Handler("GET", path, handlers)
func (rg *RouterGroup) GET(path string, handlers ...Handler) {
	rg.Handle("GET", path, handlers...)
}

// POST is a shortcut for geneva.Handler("POST", path, handlers)
func (rg *RouterGroup) POST(path string, handlers ...Handler) {
	rg.Handle("POST", path, handlers...)
}

// PUT is a shortcut for geneva.Handler("PUT", path, handlers)
func (rg *RouterGroup) PUT(path string, handlers ...Handler) {
	rg.Handle("PUT", path, handlers...)
}

// DELETE is a shortcut for geneva.Handler("DELETE", path, handlers)
func (rg *RouterGroup) DELETE(path string, handlers ...Handler) {
	rg.Handle("DELETE", path, handlers...)
}

// OPTIONS is a shortcut for geneva.Handler("OPTIONS", path, handlers)
func (rg *RouterGroup) OPTIONS(path string, handlers ...Handler) {
	rg.Handle("OPTIONS", path, handlers...)
}

// HEAD is a shortcut for geneva.Handler("HEAD", path, handlers)
func (rg *RouterGroup) HEAD(path string, handlers ...Handler) {
	rg.Handle("HEAD", path, handlers...)
}

// PATCH is a shortcut for geneva.Handler("PATCH", path, handlers)
func (rg *RouterGroup) PATCH(path string, handlers ...Handler) {
	rg.Handle("PATCH", path, handlers...)
}

// Handle registers a new request handle and middleware with the given path and method.
func (rg *RouterGroup) Handle(method, relativePath string, handlers ...Handler) {
	joined := joinPaths(rg.path, relativePath)
	combined := combineHandlers(rg.handlers, handlers)
	rg.engine.router.Handle(
		method,
		joined,
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			c := &Context{
				Writer:   w,
				Request:  r,
				Params:   ps,
				handlers: combined,
				index:    -1,
			}
			c.Next()
		},
	)
}

// Use adds middleware to the group.
func (rg *RouterGroup) Use(handlers ...Handler) {
	rg.handlers = combineHandlers(rg.handlers, handlers)
}

// Group creates a new router group.
// You should add all the routes that have common middlewares or the same path prefix.
func (rg *RouterGroup) Group(relativePath string, handlers ...Handler) *RouterGroup {
	joined := joinPaths(rg.path, relativePath)
	combined := combineHandlers(rg.handlers, handlers)
	return &RouterGroup{
		path:     joined,
		handlers: combined,
		engine:   rg.engine,
	}
}

func joinPaths(absolutePath, relativePath string) string {
	path := path.Join(absolutePath, relativePath)
	return path
}

func combineHandlers(handlers1, handlers2 Handlers) Handlers {
	s := len(handlers1) + len(handlers2)
	h := make(Handlers, 0, s)
	h = append(h, handlers1...)
	h = append(h, handlers2...)
	return h
}

/////////////
// Context //
/////////////

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort prevents pending handlers from being called.
// Note that this will not stop the current handler.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
func (c *Context) AbortWithStatus(code int) {
	c.Writer.WriteHeader(code)
	c.Abort()
}

// Set is used to store a new key/value pair exclusively for this context.
func (c *Context) Set(key string, value any) {
	if c.keys == nil {
		c.keys = make(map[string]any)
	}
	c.keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
func (c *Context) Get(key string) (value any, exists bool) {
	value, exists = c.keys[key]
	return
}

// Param returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
func (c *Context) GetQuery(key string) (string, bool) {
	if c.Request.URL.Query().Has(key) {
		return c.Request.URL.Query().Get(key), true
	}
	return "", false
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if v, ok := c.GetQuery(key); ok {
		return v
	}
	return defaultValue
}

// Query is shortcut for `c.Request.URL.Query().Get(key)`
func (c *Context) Query(key string) string {
	v, _ := c.GetQuery(key)
	return v
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
func (c *Context) GetPostForm(key string) (string, bool) {
	c.Request.ParseForm()
	if c.Request.Form.Has(key) {
		return c.Request.Form.Get(key), true
	}
	return "", false
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if v, ok := c.GetPostForm(key); ok {
		return v
	}
	return defaultValue
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string {
	v, _ := c.GetPostForm(key)
	return v
}

// String writes the given string into the response body.
func (c *Context) String(code int, msg string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}

// JSON serializes the given struct as JSON into the response body.
func (c *Context) JSON(code int, obj any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(obj)
}

////////////////
// MIDDLEWARE //
////////////////

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() Handler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
