package nanoserve

import (
	"fmt"
	"net/http"
)

type NanoServe struct {
	router *TrieRouter
}

func New() *NanoServe {
	n := &NanoServe{
		router: NewTrieRouter(10),
	}
	return n
}

func (n *NanoServe) GET(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodGet, path, h...)
}

func (n *NanoServe) POST(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodPost, path, h...)
}

func (n *NanoServe) PUT(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodPut, path, h...)
}

func (n *NanoServe) PATCH(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodPatch, path, h...)
}

func (n *NanoServe) DELETE(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodDelete, path, h...)
}

func (n *NanoServe) HEAD(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodHead, path, h...)
}

func (n *NanoServe) OPTIONS(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodOptions, path, h...)
}

func (n *NanoServe) CONNECT(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodConnect, path, h...)
}

func (n *NanoServe) TRACE(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodTrace, path, h...)
}

func (n *NanoServe) Handle(method, path string, h ...http.HandlerFunc) {
	n.addRoute(method, path, h...)
}


func (n *NanoServe) addRoute(method string, path string, handlers ...http.HandlerFunc) {
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

func (n *NanoServe) Use(pathOrHandler any, handlers ...http.HandlerFunc) {
	switch v := pathOrHandler.(type) {
	case string:
		n.router.AddMiddleware(v,handlers...)
	case http.HandlerFunc:
		all := append([]http.HandlerFunc{v}, handlers...)
		n.router.AddMiddleware("/",all...)
	}
}

// Our Main Handler which will handle the incoming request
func (n *NanoServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedHandlers := n.router.Search(r.Method, r.URL.Path)
	if len(matchedHandlers.Handler) > 0 {
		for _, f := range matchedHandlers.Handler {
			f(w, r)
		}
	}
	fmt.Fprint(w, "Not found")
}
