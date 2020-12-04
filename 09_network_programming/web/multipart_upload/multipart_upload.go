// go run multipart_upload.go
// open a web browser to http://localhost:8080/upload and upload a file

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	param    = "file"
	endpoint = "/upload"
	content  = `<html><body>` +
		`<form enctype="multipart/form-data" action="%s" method="POST">` +
		`<input type="file" name="%s"/><input type="submit" value="Upload"/>` +
		`</form></html></body>`
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, content, endpoint, param)
			return
		} else if r.Method == "POST" {
			path, err := upload(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Uploaded to %s", path)
		} else {
			http.NotFound(w, r)
			return
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func upload(r *http.Request) (string, error) {
	f, h, err := r.FormFile(param)
	if err != nil {
		return "", err
	}
	defer f.Close()

	p := filepath.Join(os.TempDir(), h.Filename)
	fw, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	if _, err = io.Copy(fw, f); err != nil {
		return "", err
	}

	return p, nil
}
