// go run runes.go

package main

import "fmt"

func main() {
	const a = `⌘`

	fmt.Printf("plain string: %s\n", a)
	fmt.Printf("quoted string: %q", a)

	fmt.Print("\nhex bytes: ")
	for i := 0; i < len(a); i++ {
		fmt.Printf("%x ", a[i])
	}
	fmt.Println()
}
