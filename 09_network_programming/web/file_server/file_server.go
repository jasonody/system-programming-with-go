// Run file server, specifying "./files" as directory to be served
// go run file_server.go ./files
// Request directories/files
// http :8080
// http :8080/a/
// http :8080/a/1
// http :8080/b/2

package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify a directory")
	}

	s, err := os.Stat(os.Args[1])
	if err == nil && !s.IsDir() {
		err = errors.New("not a directory")
	}
	if err != nil {
		log.Fatalln("Invalid path:", err)
	}

	http.Handle("/", http.FileServer(http.Dir(os.Args[1])))
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
