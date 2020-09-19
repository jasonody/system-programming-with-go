// go run pseudo_term.go stack.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jasonody/system-programming-with-go/06_pseudo_terminals/pseudo_term/command"
)

func init() {
	command.Register(&Stack{})
}

type Stack struct {
	data []string
}

func (s *Stack) push(values ...string) {
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

func (s *Stack) peek() (string, bool) {
	if len(s.data) == 0 {
		return "", false
	}

	v := s.data[len(s.data)-1]
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
	case "pop", "peek", "length":
		return len(args) == 0
	case "push":
		return len(args) > 0
	default:
		return false
	}
}

func (s *Stack) Run(r io.Reader, w io.Writer, args []string) (exit bool) {
	if l := len(args); l < 2 || !s.isValid(args[1], args[2:]) {
		fmt.Fprintf(w, "Use `stack push <something>`, `stack pop`, `stack peek`, or `stack length`\n")
		return false
	}

	switch args[1] {
	case "length":
		{
			fmt.Fprintf(w, "Length: %d", s.length())
			break
		}
	case "push":
		{
			s.push(args[2:]...)
			break
		}
	case "peek":
		{
			if v, ok := s.peek(); !ok {
				fmt.Fprintf(w, "Stack is empty!")
			} else {
				fmt.Fprintf(w, "Retrieved: `%s`\n", v)
			}
			break
		}
	case "pop":
		{
			if v, ok := s.pop(); !ok {
				fmt.Fprintf(w, "Stack is empty!")
			} else {
				fmt.Fprintf(w, "Retrieved: `%s`\n", v)
			}
			break
		}
	}
	return false
}

// get the path to the signed in user's home directory
func (s *Stack) getPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".pseudotermstack"), nil
}

func (s *Stack) Startup(w io.Writer) error {
	path, err := s.getPath()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	s.data = s.data[:0] // truncate data
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s.push(string(scanner.Bytes()))
	}

	return nil
}

func (s *Stack) Shutdown(w io.Writer) error {
	path, err := s.getPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) // O_TRUNC will truncate file
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range s.data {
		if _, err := fmt.Fprintln(f, v); err != nil {
			return err
		}
	}

	return nil
}
