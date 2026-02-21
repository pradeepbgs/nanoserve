package nanoserve

import (
	"net/http"
)

type HandlerFuntcion func(http.ResponseWriter, *http.Request, func())

type NanoServe struct {
	router *TrieRouter
}

func New() *NanoServe {
	n := &NanoServe{
		router: NewTrieRouter(),
	}
	return n
}

func (n *NanoServe) GET(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodGet, path, h...)
}

func (n *NanoServe) POST(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodPost, path, h...)
}

func (n *NanoServe) PUT(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodPut, path, h...)
}

func (n *NanoServe) PATCH(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodPatch, path, h...)
}

func (n *NanoServe) DELETE(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodDelete, path, h...)
}

func (n *NanoServe) HEAD(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodHead, path, h...)
}

func (n *NanoServe) OPTIONS(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodOptions, path, h...)
}

func (n *NanoServe) CONNECT(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodConnect, path, h...)
}

func (n *NanoServe) TRACE(path string, h ...HandlerFuntcion) {
	n.addRoute(http.MethodTrace, path, h...)
}

func (n *NanoServe) Handle(method, path string, h ...HandlerFuntcion) {
	n.addRoute(method, path, h...)
}

func (n *NanoServe) addRoute(method string, path string, handlers ...HandlerFuntcion) {
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

func (n *NanoServe) Use(pathOrHandler any, handlers ...HandlerFuntcion) {
	switch v := pathOrHandler.(type) {
	case string:
		n.router.AddMiddleware(v, handlers...)
	case HandlerFuntcion:
		all := append([]HandlerFuntcion{v}, handlers...)
		n.router.AddMiddleware("/", all...)
	}
}

// Our Main Handler which will handle the incoming request
func (n *NanoServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	match := n.router.Search(r.Method, r.URL.Path)
	if len(match.Handler) == 0 {
		http.NotFound(w, r)
		return
	}

	handlers := match.Handler

	index := 0

	var next func()

	if len(handlers) == 1 {
		handlers[0](w, r, next)
		return
	}

	next = func() {
		if index >= len(handlers) {
			return
		}
		h := handlers[index]
		index++
		h(w, r, next)
	}
	next()
}
