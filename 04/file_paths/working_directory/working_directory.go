package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("starting directory:", wd)

	if err := os.Chdir("/"); err != nil {
		fmt.Println(err)
		return
	}

	if wd, err = os.Getwd(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wd)
}
