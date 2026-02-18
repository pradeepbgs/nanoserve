package main

import (
	"fmt"
	"github.com/pradeepbgs/nanoserve/nanoserve"
)



func main() {
	router := nanoserve.NewTrieRouter(1)

	router.AddMiddleware("/", func() {
		fmt.Println("Global Middleware")
	})
	router.AddMiddleware("/user/*", func() {
		fmt.Println("?user middleware")
	})
	router.AddRoute("GET", "/user/:id", func() {
		fmt.Println("GET /user/:id")
	})
	router.AddMiddleware("/register", func() {
		fmt.Println("/register middleware")
	})
	router.AddRoute("POST", "/register", func() {
		fmt.Println("POST register")
	})

	handler := router.MatchRoute("POST", "/register")
	if handler == nil {
		fmt.Println("Route Not Found")
		return
	}

	for _, f := range handler.Handler {
		if f != nil {
			f()
		}
	}
}
