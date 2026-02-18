package nanoserve

import (
	"strings"
)

type RouteMatch struct {
	Handler []func()
}

type Node struct {
	children    map[string]*Node
	isEndOfWord bool
	handlers    map[string]func()
	middlewares []func()
}

func (n *Node) NewNode() *Node {
	return &Node{
		children:    make(map[string]*Node),
		handlers:    make(map[string]func()),
		middlewares: make([]func(), 0, 10),
		isEndOfWord: false,
	}
}

type Trie struct {
	root *Node
}

func NewTrieRouter(middlewares_size int) *Trie {
	return &Trie{
		root: &Node{
			children:    map[string]*Node{},
			handlers:    map[string]func(){},
			isEndOfWord: false,
			middlewares: make([]func(), 0, middlewares_size),
		},
	}
}

func (r *Trie) AddMiddleware(path string, handlers ...func()) {
	node := r.root
	if path == "/" {
		node.middlewares = append(node.middlewares, handlers...)
		return
	}

	pathSegments := strings.SplitSeq(path, "/")

	for element := range pathSegments {
		if element == "" {
			break
		}
		key := element

		if strings.HasPrefix(element, ":") {
			key = ":"
		} else if strings.HasPrefix(element, "*") {
			node.middlewares = append(node.middlewares, handlers...)
		}

		if node.children[key] == nil {
			node.children[key] = node.NewNode()
		}
		node = node.children[key]
	}
	node.middlewares = append(node.middlewares, handlers...)
}

func (r *Trie) AddRoute(method string, path string, handler func()) {
	node := r.root

	var pathSegments []string = strings.Split(path, "/")

	if path == "/" {
		node.isEndOfWord = true
		node.handlers[method] = handler
		return
	}

	for _, element := range pathSegments {
		if element == "" {
			continue
		}
		key := element

		if strings.HasPrefix(element, ":") {
			key = ":"
		}

		if node.children[key] == nil {
			node.children[key] = node.NewNode()
		}

		node = node.children[key]
	}
	node.isEndOfWord = true
	node.handlers[method] = handler
}

func (r *Trie) MatchRoute(method string, path string) *RouteMatch {
	node := r.root

	pathSegments := strings.Split(path, "/")
	collectedHandlers := append([]func(){}, node.middlewares...)

	for _, element := range pathSegments {
		if element == "" {
			continue
		}
		if node.children[element] != nil {
			node = node.children[element]
		} else if node.children[":"] != nil {
			node = node.children[":"]
		} else if node.children["*"] != nil {
			break
		} else {
			return &RouteMatch{Handler: collectedHandlers}
		}
	}

	methodHandler := node.handlers[method]
	if methodHandler != nil {
		collectedHandlers = append(collectedHandlers, methodHandler)
		return &RouteMatch{Handler: collectedHandlers}
	}
	return nil
}
