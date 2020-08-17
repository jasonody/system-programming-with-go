// go run resetting_buffer_and_strings.go

package main

import (
	"bytes"
	"fmt"
)

func main() {
	b := bytes.NewBuffer(nil)
	b.WriteString("One")
	s1 := b.String()
	b.WriteString("Two")
	s2 := b.String()
	b.Reset()
	b.WriteString("Hey!") // does not change s1 or s2 as strings are immutable (the bytes were copied from the buffer)
	s3 := b.String()

	fmt.Println(s1, s2, s3) // prints "One OneTwo Hey!"
}
