// go build goroutines.go && ./goroutines
// Output:
// main start
// go start

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	go func() {
		defer fmt.Println("go end (deferred)") // deferred doesn't get executed
		fmt.Println("go start")
		os.Exit(1) // all goroutines, including the main one, will terminate immediately
	}()

	defer fmt.Println("main end (deferred") // deferred doesn't get executed
	fmt.Println("main start")
	time.Sleep(time.Second)
	fmt.Println("main end") // program already exited
}
