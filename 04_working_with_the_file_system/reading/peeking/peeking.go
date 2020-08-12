// go run peeking.go ../sample.txt

package main

import (
	"bufio"
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
	defer file.Close() // ensure close to avoid leaks

	reader := bufio.NewReader(file)
	var rowCount int

	for err == nil {
		var b []byte
		for more := true; err == nil && more; {
			b, more, err = reader.ReadLine()
			if err == nil {
				fmt.Println(string(b))
			}
		}
		// each time more is false, a line has been completely read
		if err == nil {
			fmt.Println()
			rowCount++
		}
	}

	if err != nil && err != io.EOF {
		fmt.Println("\nError:", err)
		return
	}

	fmt.Println("\nRow count:", rowCount)
}
