package command

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrCommandNotFound is returned when a given command name isn't found
	ErrCommandNotFound = errors.New("Command not found")
)

// Set represents a set of commands that can be executed
type Set struct {
	commands map[string]Command
}

// NewSet returns a set of the given commands
func NewSet(commands []Command) (*Set, error) {
	commandsByName := make(map[string]Command)
	for _, c := range commands {
		name := strings.ToLower(c.Name())

		if name == "" {
			return nil, errors.New("Name must not be blank")
		}

		if _, ok := commandsByName[name]; ok {
			return nil, fmt.Errorf("Command with name %s provided twice", name)
		}

		commandsByName[name] = c
	}

	return &Set{
		commands: commandsByName,
	}, nil
}

// Get returns the command with the given name, or nil if no command with the given name exists
func (s *Set) Get(name string) Command {
	return s.commands[strings.ToLower(name)]
}

// Execute executes a command in the set. It assumes that the first arg is the command name and
// that the rest of the args should be provided to the command
func (s *Set) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No command name provided in args")
	}

	name := args[0]
	cmd := s.Get(name)
	if cmd == nil {
		return ErrCommandNotFound
	}

	rest := args[1:]
	return cmd.Execute(ctx, rest)
}
