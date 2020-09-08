// go run pseudo_term.go

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

	"github.com/agnivade/levenshtein"
)

func exit(w io.Writer, args []string) bool {
	fmt.Fprintf(w, "Goodbye! :)\n")
	return true
}

func help(w io.Writer, args []string) bool {
	fmt.Fprintln(w, "Available commands:")
	for _, cmd := range cmds {
		fmt.Fprintf(w, "  - %-15s %s\n", cmd.Name, cmd.Help)
	}
	return false
}

func shuffle(w io.Writer, args []string) bool {
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
	})
	for i := range args {
		if i > 0 {
			fmt.Fprint(w, " ") // add a space to seperate words
		}
		fmt.Fprintf(w, "%s", args[i])
	}
	fmt.Fprintln(w)
	return false
}

func print(w io.Writer, args []string) bool {
	if len(args) != 1 {
		fmt.Fprintln(w, "Please specify one file.")
		return false
	}

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(w, "Cannot open %s: %s\n", args[0], err)
		return false
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		fmt.Fprintf(w, "Cannot print %s: %s\n", args[0], err)
		return false
	}

	fmt.Fprintln(w)
	return false
}

type cmd struct {
	Name   string                                // the command name
	Help   string                                // a description of the command
	Action func(w io.Writer, args []string) bool // command handler
}

func (c cmd) Match(s string) bool {
	return c.Name == s
}

func (c cmd) Run(w io.Writer, args []string) bool {
	return c.Action(w, args)
}

var cmds = make([]cmd, 0, 10)

// init() gets called before main()
func init() {
	cmds = append(cmds,
		cmd{
			Name:   "exit",
			Help:   "Exits the application",
			Action: exit,
		},
		cmd{
			Name:   "help",
			Help:   "Shows available commands",
			Action: help,
		},
		cmd{
			Name:   "shuffle",
			Help:   "Shuffles a list of strings",
			Action: shuffle,
		},
		cmd{
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

	fmt.Fprint(w, "Welcome to PseudoTerm!\n")

	for {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get working directory", err)
			return
		}
		fmt.Fprintf(w, "\n[%s] >", filepath.Base(pwd))

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

		idx := -1
		for i := range cmds {
			if !cmds[i].Match(args[0]) {
				continue
			}
			idx = i
			break
		}

		if idx == -1 {
			commandNotFound(w, args[0])
			continue
		}

		if cmds[idx].Run(w, args[1:]) { // execute command and exit it returns true
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

func commandNotFound(w io.Writer, cmd string) {
	var suggestions []string

	for _, c := range cmds {
		if levenshtein.ComputeDistance(c.Name, cmd) < 3 {
			suggestions = append(suggestions, c.Name)
		}
	}

	fmt.Fprintf(w, "Command %q not found. ", cmd)
	if len(suggestions) == 0 {
		fmt.Fprint(w, "Use `help` for available commands.\n")
		return
	}
	fmt.Fprint(w, "Maybe you meant: ")
	for i, suggestion := range suggestions {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprintf(w, "%s", suggestion)
	}
	fmt.Fprintln(w)
}
