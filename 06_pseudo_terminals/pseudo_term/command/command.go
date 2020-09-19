package command

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/agnivade/levenshtein"
)

// command represents a terminal command
type Command interface {
	Startup(w io.Writer) error
	Shutdown(w io.Writer) error
	GetName() string
	GetHelp() string
	Run(input io.Reader, output io.Writer, args []string) (exit bool)
}

// ErrDuplicateCommand is returned when two commands have the same name
var ErrDuplicateCommand = errors.New("Duplicate command")

var commands []Command

// Startup execute startup for all commands
func Startup(w io.Writer) {
	for _, c := range commands {
		if err := c.Startup(w); err != nil {
			fmt.Fprintf(w, "%s: startup error: %s", c.GetName(), err)
		}
	}
}

// Shutdown executes shutdown for all commands
func Shutdown(w io.Writer) {
	for _, c := range commands {
		if err := c.Shutdown(w); err != nil {
			fmt.Fprintf(w, "%s: shutdown error: %s", c.GetName(), err)
		}
	}
}

// Register adds the Command to the command list
func Register(command Command) error {
	name := command.GetName()
	for i, c := range commands {
		// unique commands in alphabetical order
		switch strings.Compare(c.GetName(), name) {
		case 0:
			return ErrDuplicateCommand
		case 1:
			commands = append(commands, nil)
			copy(commands[i+1:], commands[i:])
			commands[i] = command
			return nil
		case -1:
			continue
		}
	}
	commands = append(commands, command)
	return nil
}

// Base is a basic Command that runs a closure
type Base struct {
	Name, Help string
	Action     func(input io.Reader, output io.Writer, args []string) (exit bool)
}

// Startup does nothing
func (b Base) Startup(w io.Writer) error { return nil }

// Shutdown does nothing
func (b Base) Shutdown(w io.Writer) error { return nil }

func (b Base) String() string { return b.Name }

// GetName returns the Name
func (b Base) GetName() string { return b.Name }

// GetHelp returns the Help
func (b Base) GetHelp() string { return b.Help }

// Run calls the closure
func (b Base) Run(input io.Reader, output io.Writer, args []string) (exit bool) {
	return b.Action(input, output, args)
}

// GetCommand returns the command with the given name
func GetCommand(name string) Command {
	for _, c := range commands {
		if c.GetName() == name {
			return c
		}
	}

	// if no match, return suggest command
	return suggest
}

var suggest = Base{
	Action: func(in io.Reader, w io.Writer, args []string) bool {
		var suggestions []string

		for _, c := range commands {
			name := c.GetName()
			if levenshtein.ComputeDistance(name, args[0]) < 3 {
				suggestions = append(suggestions, name)
			}
		}

		fmt.Fprintf(w, "Command %q not found. ", args[0])
		if len(suggestions) == 0 {
			fmt.Fprint(w, "Use `help` for available commands.\n")
			return false
		}
		fmt.Fprint(w, "Maybe you meant: ")
		for i, suggestion := range suggestions {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprintf(w, "%s", suggestion)
		}
		fmt.Fprintln(w)
		return false
	},
}

func init() {
	Register(Base{Name: "help", Help: "Shows available commands", Action: helpAction})
	Register(Base{Name: "exit", Help: "Exits the application", Action: exitAction})
}

func helpAction(in io.Reader, w io.Writer, args []string) bool {
	fmt.Fprintln(w, "Available commands:")
	for _, cmd := range commands {
		fmt.Fprintf(w, "  - %-15s %s\n", cmd.GetHelp(), cmd.GetHelp())
	}
	return false
}

func exitAction(in io.Reader, w io.Writer, args []string) bool {
	fmt.Fprintf(w, "Goodbye! :)\n")
	return true
}
