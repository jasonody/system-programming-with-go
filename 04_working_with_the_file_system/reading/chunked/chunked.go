// go run chunked.go ../sample.txt

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please specify a file path.")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Cannot open:", err)
		return
	}
	defer f.Close() // ensure close to avoid leaks

	// var (
	// 	b = make([]byte, 16)
	// )
	b := make([]byte, 16) // same as above commented out code
	for n := 0; err == nil; {
		n, err = f.Read(b)
		if err == nil {
			// fmt.Println("New chunk...", n)
			fmt.Print(string(b[:n])) // only print what's been read
		}
	}

	if err != nil && err != io.EOF { // we expect EOF at some point
		fmt.Println("\n\n Error:", err)
	} else {
		fmt.Print("\n___________\nEND OF FILE\n")
	}
}
