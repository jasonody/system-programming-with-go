// go run seeking.go ../sample.txt
// Should print "Read data: Some"

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

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Cannot open:", err)
		return
	}
	defer file.Close()

	file.Seek(9, io.SeekStart)
	data := make([]byte, 4)
	_, err = file.Read(data)

	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Println("Read data:", string(data))
}
