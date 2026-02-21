package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pradeepbgs/nanoserve"
)

func Logger(w http.ResponseWriter, r *http.Request, next func()) {
	fmt.Println("starting path ", r.URL.Path)
	start := time.Now()
	fmt.Println("user calling next")
	next()
	fmt.Println("Request took:", time.Since(start))
}

func main() {
	router := nanoserve.New()

	router.GET("/user", Logger, func(w http.ResponseWriter, r *http.Request, next func()) {
		fmt.Fprint(w, "Hello from /")
	})

	router.Run(":3000")
}
