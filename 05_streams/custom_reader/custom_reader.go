// go run custom_reader.go

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"
)

func NewAngryReader(r io.Reader) *AngryReader {
	return &AngryReader{r: r}
}

type AngryReader struct {
	r io.Reader
}

func (a *AngryReader) Read(b []byte) (int, error) {
	n, err := a.r.Read(b)
	for runeChar, i, width := rune(0), 0, 0; i < n; i += width {
		runeChar, width = utf8.DecodeRune(b[i:]) // [i:] = from i to end of slice
		fmt.Printf("Rune: %c, width: %d\n", runeChar, width)
		if !unicode.IsLetter(runeChar) {
			continue
		}

		runeUpper := unicode.ToUpper(runeChar)
		if widthUpper := utf8.EncodeRune(b[i:], runeUpper); width != widthUpper {
			return n, fmt.Errorf("%c->%c, size mismatch %d->%d", runeChar, runeUpper, width, widthUpper)
		}
	}

	return n, err
}

func main() {
	a := NewAngryReader(strings.NewReader("⌘Hello, playground!⌘"))
	b, err := ioutil.ReadAll(a)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(b))
}
