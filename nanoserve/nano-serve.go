package nanoserve

import (
	"fmt"
	"net/http"
)

type NanoServe struct {
	router *TrieRouter
}

var HttpMethods = []string{"GET", "POST", "PUT", "DELETE"}

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

func (n *NanoServe) DELETE(path string, h ...http.HandlerFunc) {
	n.addRoute(http.MethodDelete, path, h...)
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
