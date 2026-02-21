package nanoserve

import (
	"net/http"
	"strings"
)

type RouteMatch struct {
	Handler []http.HandlerFunc
}

type Node struct {
	children    map[string]*Node
	isEndOfWord bool
	handlers    map[string]http.HandlerFunc
	middlewares []http.HandlerFunc
}

func newNode() *Node {
	return &Node{
		children:    make(map[string]*Node),
		handlers:    make(map[string]http.HandlerFunc),
		middlewares: make([]http.HandlerFunc, 0, 10),
	}
}

type TrieRouter struct {
	root              *Node
	globalMiddlewares []http.HandlerFunc
}

func NewTrieRouter(middlewaresSize int) *TrieRouter {
	return &TrieRouter{
		root: &Node{
			children:    make(map[string]*Node),
			handlers:    make(map[string]http.HandlerFunc),
			middlewares: make([]http.HandlerFunc, 0, middlewaresSize),
		},
	}
}

func (r *TrieRouter) AddMiddleware(path string, handlers ...http.HandlerFunc) {
	node := r.root

	if path == "/" {
		r.globalMiddlewares = append(r.globalMiddlewares, handlers...)
		return
	}

	segments := strings.Split(path, "/")

	for _, element := range segments {
		if element == "" {
			continue
		}

		key := element
		if strings.HasPrefix(element, ":") {
			key = ":"
		} else if strings.HasPrefix(element, "*") {
			node.middlewares = append(node.middlewares, handlers...)
		}

		if node.children[key] == nil {
			node.children[key] = newNode()
		}
		node = node.children[key]
	}

	node.middlewares = append(node.middlewares, handlers...)
}

func (r *TrieRouter) Insert(method string, path string, handler http.HandlerFunc) {
	node := r.root

	if path == "/" {
		node.isEndOfWord = true
		node.handlers[method] = handler
		return
	}

	segments := strings.Split(path, "/")

	for _, element := range segments {
		if element == "" {
			continue
		}

		key := element
		if strings.HasPrefix(element, ":") {
			key = ":"
		}

		if node.children[key] == nil {
			node.children[key] = newNode()
		}

		node = node.children[key]
	}

	node.isEndOfWord = true
	node.handlers[method] = handler
}

func (r *TrieRouter) Search(method string, path string) *RouteMatch {
	node := r.root

	segments := strings.Split(path, "/")

	collected := append([]http.HandlerFunc{}, r.globalMiddlewares...)

	for _, element := range segments {
		if element == "" {
			continue
		}

		if node.children[element] != nil {
			node = node.children[element]
		} else if node.children[":"] != nil {
			node = node.children[":"]
		} else if node.children["*"] != nil {
			node = node.children["*"]
			break
		} else {
			return &RouteMatch{Handler: collected}
		}

		collected = append(collected, node.middlewares...)
	}

	if h := node.handlers[method]; h != nil {
		collected = append(collected, h)
		return &RouteMatch{Handler: collected}
	}

	return nil
}