package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please supply a path and some contents.")
		return
	}

	// file contents needs to be casted to a byte slice
	if err := ioutil.WriteFile(os.Args[1], []byte(os.Args[2]), 0644); err != nil {
		fmt.Println("Error writing", err)
	}
}
