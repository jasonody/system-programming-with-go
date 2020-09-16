// go run pseudo_term.go stack.go

package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/jasonody/system-programming-with-go/06_pseudo_terminals/pseudo_term/command"
)

func init() {
	command.Register(&Stack{})
}

type Stack struct {
	data []string
}

func (s *Stack) push(values []string) {
	// s.data = append(s.data, values...) // original stored each string in array/value separately: ["this" "is" "not" "one" "string"]
	s.data = append(s.data, strings.Join(values, " "))
}

func (s *Stack) pop() (string, bool) {
	if len(s.data) == 0 {
		return "", false
	}

	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func (s *Stack) length() int {
	return len(s.data)
}

func (s *Stack) GetName() string {
	return "stack"
}

func (s *Stack) GetHelp() string {
	return "stack-like memory storage"
}

func (s *Stack) isValid(cmd string, args []string) bool {
	switch cmd {
	case "pop", "length":
		return len(args) == 0
	case "push":
		return len(args) > 0
	default:
		return false
	}
}

func (s *Stack) Run(r io.Reader, w io.Writer, args []string) (exit bool) {
	if l := len(args); l < 2 || !s.isValid(args[1], args[2:]) {
		fmt.Fprintf(w, "Use `stack push <something>` or `stack pop`\n")
		return false
	}

	if args[1] == "length" {
		fmt.Fprintf(w, "Length: %d", s.length())
		return false
	}

	if args[1] == "push" {
		s.push(args[2:])
		return false
	}

	if v, ok := s.pop(); !ok {
		fmt.Fprintf(w, "Stack is empty!")
	} else {
		fmt.Fprintf(w, "Retrieved: `%s`\n", v)
	}
	return false
}
