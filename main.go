package main

import (
	"fmt"
	"net/http"

	"github.com/pradeepbgs/nanoserve/nanoserve"
)



func main() {
	router := nanoserve.New()

	router.GET("/",func (w http.ResponseWriter, r *http.Request)  {
		fmt.Fprint(w, "Hello from /")
	})

	router.Run(":3000")
}
