// go run pipe.go

package main

import (
	"fmt"
	"io"
)

func main() {
	preader, pwriter := io.Pipe()

	go func(writer io.WriteCloser) {
		defer writer.Close()
		for _, s := range []string{"a string", "another string", "last one"} {
			fmt.Printf("-> writing %q\n", s)
			fmt.Fprint(writer, s)
		}
	}(pwriter)

	var err error
	for length, buffer := 0, make([]byte, 100); err == nil; {
		fmt.Println("<- waiting...")
		length, err = preader.Read(buffer) // read from pipe into buffer
		if err == nil {
			fmt.Printf("<- received %q\n", string(buffer[:length]))
		}
	}

	if err != nil && err != io.EOF {
		fmt.Println("error:", err)
	}
}
