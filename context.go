package nanoserve

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	params map[string]string

	handlers []HandlerFunction
	index    int

	contextData map[string]interface{}
}

func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		return
	}
	c.handlers[c.index](c)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) Set(key string, value any) {
	c.contextData[key] = value
}

func (c *Context) Get(key string) interface{} {
	return c.contextData[key]
}

func (c *Context) Text(text string, code int) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(text))
}

func (c *Context) Url() *url.URL {
	return c.Request.URL
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	val := c.params[key]
	if val != "" {
		return val
	}
	return ""
}

func (c *Context) Json(data any, code int) {
	c.Writer.WriteHeader(code)
	err := json.NewEncoder(c.Writer).Encode(data)
	if err != nil {
		panic("error during encoding json")
	}
}
