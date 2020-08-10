// go run all_at_once.go ../sample.txt

// Not a good idea to do this is the file is large or it's size is unknown

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please specify a file path.")
		return
	}

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return
	}

	fmt.Println(string(b))
}
