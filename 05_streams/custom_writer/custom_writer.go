package main

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"unicode"
	"unicode/utf8"
)

func NewScrambleWriter(w io.Writer, r *rand.Rand, c float64) *ScrambleWriter {
	return &ScrambleWriter{writer: w, random: r, chance: c}
}

type ScrambleWriter struct {
	writer io.Writer
	random *rand.Rand
	chance float64
}

func (s *ScrambleWriter) scrambleWrite(runes []rune, wordSeparator rune) (n int, err error) {
	// scramble after first letter
	for i := 1; i < len(runes)-1; i++ {
		if s.random.Float64() > s.chance {
			continue
		}
		j := s.random.Intn(len(runes)-1) + 1
		runes[i], runes[j] = runes[j], runes[i]
	}

	if wordSeparator != 0 {
		runes = append(runes, wordSeparator)
	}

	var b = make([]byte, 10)
	for _, r := range runes {
		v, err := s.writer.Write(b[:utf8.EncodeRune(b, r)]) // encode the rune into b and then only write the bytes the rune occupies to the writer
		// fmt.Println(string(b))
		if err != nil {
			return n, err
		}

		n += v
	}
	return
}

func (s *ScrambleWriter) Write(b []byte) (n int, err error) {
	var runes = make([]rune, 0, 10) // initial capacity (10) doen't have to be accurate as "append" will grow slice when needed
	for readRune, i, readRuneSize := rune(0), 0, 0; i < len(b); i += readRuneSize {
		readRune, readRuneSize = utf8.DecodeRune(b[i:])
		if unicode.IsLetter(readRune) {
			runes = append(runes, readRune)
			continue
		}

		// rune is not a letter
		v, err := s.scrambleWrite(runes, readRune)
		if err != nil {
			return n, err
		}

		n += v
		runes = runes[:0] // empty runes slice
	}

	// handle scenario where last rune is a letter
	if len(runes) != 0 {
		v, err := s.scrambleWrite(runes, 0)
		if err != nil {
			return n, err
		}

		n += v
	}
	return
}

func main() {
	var s strings.Builder
	w := NewScrambleWriter(&s, rand.New(rand.NewSource(1)), 0.5)
	fmt.Fprint(w, "Hello! This is scrambled text. It was scrambled using the ScrableWriter")

	fmt.Println(s.String())
}
