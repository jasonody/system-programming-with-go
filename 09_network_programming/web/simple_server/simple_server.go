// Run web-server:
// go run simple_server.go
// Make requests:
// http :8080/hello
// http :8080/bye
// http :8080/error
// http :8080/custom

package main

import (
	"fmt"
	"log"
	"net/http"
)

type customHandler int

func (c *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", *c)
	*c++
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	mux.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bye!")
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader((http.StatusInternalServerError))
		fmt.Fprintf(w, "An error occured.")
	})
	mux.Handle("/custom", new(customHandler))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
