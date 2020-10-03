// go build working_dir.go && ./working_dir

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Working directory: ", wd)
	fmt.Println("Application: ", filepath.Join(wd, os.Args[0]))

	// create a new directory
	d := filepath.Join(wd, "test")
	if err := os.Mkdir(d, 0755); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Created directory: ", d)

	// change current directory
	if err := os.Chdir(d); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("New working directory: ", d)
}
