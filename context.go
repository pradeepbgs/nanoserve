package nanoserve

import (
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	handlers []HandlerFunction
	index    int
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

func (c *Context) Text(text string, code int) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(text))
}