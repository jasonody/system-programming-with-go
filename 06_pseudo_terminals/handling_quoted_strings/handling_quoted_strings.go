// go run handling_quoted_strings.go

package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	s := bufio.NewScanner(strings.NewReader("hi 'this day'"))
	// s := bufio.NewScanner(strings.NewReader(`And the 'cats in
	// the craddle' and "the silver's spoon" nasa put a man on "the moon`))
	s.Split(ScanArgs) // set the split function for the scanner (default is ScanLines)

	for s.Scan() {
		fmt.Printf("%q\n", s.Text())
	}
	fmt.Println(s.Err())
}

var ErrClosingQuote = errors.New("Missing closing quote")

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

// this implementation does not handle quoted string that spans across chunks
func ScanArgs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// first space
	start, first := 0, rune(0)
	for width := 0; start < len(data); start += width {
		first, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(first) {
			break
		}
	}

	// skip quote
	if isQuote(first) {
		start++
	}

	// loop until arg end character
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		fmt.Printf("first: %q r: %q | ", string(first), string(r))
		if isFirstCharQuote := isQuote(first); !isFirstCharQuote && unicode.IsSpace(r) || isFirstCharQuote && r == first {
			return i + width, data[start:i], nil
		}
	}

	// token from EOF
	if atEOF && len(data) > start {
		if isQuote(first) {
			err = ErrClosingQuote
		}
		return len(data), data[start:], err
	}

	if isQuote(first) {
		start--
	}
	return start, nil, nil
}
