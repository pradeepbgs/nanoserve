package nanoserve

import (
	"net/http"
)

type HandlerFunction func(*Context)

type NanoServe struct {
	router *TrieRouter
}

func New() *NanoServe {
	return &NanoServe{
		router: NewTrieRouter(),
	}
}

func (n *NanoServe) GET(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodGet, path, h...)
}

func (n *NanoServe) POST(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodPost, path, h...)
}

func (n *NanoServe) PUT(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodPut, path, h...)
}

func (n *NanoServe) PATCH(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodPatch, path, h...)
}

func (n *NanoServe) DELETE(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodDelete, path, h...)
}

func (n *NanoServe) HEAD(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodHead, path, h...)
}

func (n *NanoServe) OPTIONS(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodOptions, path, h...)
}

func (n *NanoServe) CONNECT(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodConnect, path, h...)
}

func (n *NanoServe) TRACE(path string, h ...HandlerFunction) {
	n.addRoute(http.MethodTrace, path, h...)
}

func (n *NanoServe) Handle(method, path string, h ...HandlerFunction) {
	n.addRoute(method, path, h...)
}

func (n *NanoServe) addRoute(method string, path string, handlers ...HandlerFunction) {
	if len(handlers) == 0 {
		panic("route must have at least one handler")
	}

	middlewareFunctions := handlers[:len(handlers)-1]
	if len(middlewareFunctions) > 0 {
		n.router.AddMiddleware(path, middlewareFunctions...)
	}

	handler := handlers[len(handlers)-1]
	n.router.Insert(method, path, handler)
}

func (n *NanoServe) Run(addr string) error {
	return http.ListenAndServe(addr, n)
}

func (n *NanoServe) Use(pathOrHandler any, handlers ...HandlerFunction) {
	switch v := pathOrHandler.(type) {
	case string:
		n.router.AddMiddleware(v, handlers...)
	case HandlerFunction:
		all := append([]HandlerFunction{v}, handlers...)
		n.router.AddMiddleware("/", all...)
	}
}

// Our Main Handler which will handle the incoming request
func (n *NanoServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	match := n.router.Search(r.Method, r.URL.Path)
	if match == nil || len(match.Handler) == 0 {
		http.NotFound(w, r)
		return
	}

	c := &Context{
		Writer: w,
		Request: r,
		handlers: match.Handler,
		index: 0,
	}
	c.handlers[0](c)
}
