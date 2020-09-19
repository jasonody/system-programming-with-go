// go run pseudo_term.go stack.go

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/jasonody/system-programming-with-go/06_pseudo_terminals/pseudo_term/command"
)

type color int

func (c color) Start(w io.Writer) {
	fmt.Fprintf(w, "\x1b[%dm", c)
}

func (c color) End(w io.Writer) {
	fmt.Fprintf(w, "\x1b[%dm", Reset)
}

func (c color) Sprintf(w io.Writer, format string, args ...interface{}) {
	c.Start(w)
	fmt.Fprintf(w, format, args...)
	c.End(w)
}

// list of colors
const (
	Reset   color = 0
	Red     color = 31
	Green   color = 32
	Yellow  color = 33
	Blue    color = 34
	Magenta color = 35
	Cyan    color = 36
	White   color = 37
)

func shuffle(r io.Reader, w io.Writer, args []string) bool {
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
	})

	for i := range args {
		if i > 0 {
			fmt.Fprint(w, " ") // add a space to seperate words
		}
		var f func(w io.Writer, format string, args ...interface{})
		if i%2 == 0 {
			f = Red.Sprintf
		} else {
			f = Green.Sprintf
		}
		f(w, "%s", args[i])
	}
	fmt.Fprintln(w)
	return false
}

func print(r io.Reader, w io.Writer, args []string) bool {
	if len(args) != 2 {
		fmt.Fprintln(w, "Please specify one file.")
		return false
	}

	f, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(w, "Cannot open %s: %s\n", args[1], err)
		return false
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		fmt.Fprintf(w, "Cannot print %s: %s\n", args[1], err)
		return false
	}

	fmt.Fprintln(w)
	return false
}

// init() gets called before main()
func init() {
	command.Register(command.Base{
		Name:   "shuffle",
		Help:   "Shuffles a list of strings",
		Action: shuffle,
	})
	command.Register(command.Base{
		Name:   "print",
		Help:   "Prints a file",
		Action: print,
	})
}

func main() {
	s := bufio.NewScanner(os.Stdin) // bufio.Scanner is a line reader
	w := os.Stdout
	args := argsScanner{}
	b := bytes.Buffer{}

	command.Startup(w)
	defer command.Shutdown(w) // this is executed before returning

	fmt.Fprint(w, "Welcome to PseudoTerm!\n")

	for {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get working directory", err)
			return
		}

		Blue.Sprintf(w, "\n[%s]> ", filepath.Base(pwd))

		args.Reset()
		b.Reset()

		for {
			s.Scan()
			b.Write(s.Bytes())
			extra := args.Parse(&b)
			if extra == "" {
				break
			}
			b.WriteString(extra)
			fmt.Println(extra)
		}

		if command.GetCommand(args[0]).Run(os.Stdin, w, args) {
			fmt.Fprintln(w)
			return
		}
	}
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

type argsScanner []string

func (a *argsScanner) Reset() { *a = (*a)[0:0] }

func (a *argsScanner) Parse(r io.Reader) string {
	s := bufio.NewScanner(r)
	s.Split(ScanArgs)
	for s.Scan() {
		*a = append(*a, s.Text())
	}

	if err := s.Err(); err != nil {
		fmt.Println(err)
	}

	if len(*a) == 0 {
		return ""
	}

	lastArg := (*a)[len(*a)-1]
	if !isQuote(rune(lastArg[0])) {
		return ""
	}

	*a = (*a)[:len(*a)-1]
	return lastArg + "\n"
}
