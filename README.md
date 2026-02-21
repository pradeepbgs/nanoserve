# nanoServe

nanoServe is a lightweight HTTP router written in Go. It provides a trie-based routing engine with support for middleware, multiple HTTP methods, and flexible route handling.

The goal of nanoServe is to remain simple, fast, and easy to extend while maintaining clean architectural separation between routing and execution.

---

## Features

* Trie-based path matching
* Support for all standard HTTP methods
* Global and path-specific middleware
* Variadic handler support
* Implements `http.Handler` interface
* Minimal and extensible core design

---

## Installation

```bash
go get github.com/yourusername/nanoserve
```

---

## Quick Start

```go
package main

import (
	"net/http"
	"github.com/yourusername/nanoserve"
)

func main() {
	r := nanoserve.New()

	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Run(":8080")
}
```

Visit [http://localhost:8080/hello](http://localhost:8080/hello)

---

## Supported HTTP Methods

nanoServe supports all standard HTTP methods:

* GET
* POST
* PUT
* PATCH
* DELETE
* HEAD
* OPTIONS
* CONNECT
* TRACE

You can also use:

```go
r.Handle(method, path, handlers...)
r.ANY(path, handlers...)
```

---

## Middleware

### Global Middleware

```go
r.Use(func(w http.ResponseWriter, r *http.Request) {
	println("global middleware")
})
```

### Path-Specific Middleware

```go
r.Use("/api", func(w http.ResponseWriter, r *http.Request) {
	println("api middleware")
})
```

Middleware executes in the order they are registered.

---

## Route Parameters

nanoServe supports parameterized routes using `:` prefix.

Example:

```go
r.GET("/users/:id", handler)
```

Parameter extraction logic is handled by the trie router.

---

## Architecture Overview

nanoServe follows a layered architecture:

1. HTTP Layer – Implements `http.Handler`
2. Router Layer – Manages route registration
3. Trie Layer – Handles path matching
4. Middleware Layer – Collects and executes middleware

This separation ensures maintainability and extensibility.

---

## Project Structure

```
nanoserve/
├── router.go
├── trie.go
├── server.go
├── middleware.go
├── context.go (planned)
```

---

## Roadmap

* Context abstraction
* Improved middleware chaining
* Route groups
* Subrouters
* 405 Method Not Allowed support
* Automatic OPTIONS handling
* Performance benchmarking

---

## Philosophy

nanoServe aims to:

* Keep the core minimal
* Avoid unnecessary abstractions
* Maintain clean separation of concerns
* Stay idiomatic to Go

---

## License

MIT License

---

## Contributing

Contributions are welcome. Feel
