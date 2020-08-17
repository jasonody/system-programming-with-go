// go run strings_builder.go

package main

import (
	"fmt"
	"strings"
)

func main() {
	b := strings.Builder{}
	b.WriteString("One")
	s1 := b.String()
	b.WriteString("Two")
	s2 := b.String()
	b.Reset()
	b.WriteString("Hey!")
	s3 := b.String()
	fmt.Println(s1, s2, s3) // prints "One OneTwo Hey!"

	c := b
	c.WriteString("Panic!") // panic: strings: illegal use of non-zero Builder copied by value
}
