// Run the server
// go run http_methods.go
// Make some requests
// http :8080/path1
// http POST :8080/path1
// http :8080/path2
// http POST :8080/path2
// http :8080/notfound

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type methodHandlerMap map[string]http.Handler

func (m methodHandlerMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := m[strings.ToUpper(r.Method)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func main() {
	http.Handle("/path1", methodHandlerMap{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Showing record")
		}),
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Creating record")
		}),
	})

	http.Handle("/path2", methodHandlerMap{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Showing record for a different path")
		}),
		http.MethodPost: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Creating record for a different path")
		}),
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
