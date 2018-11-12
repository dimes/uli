package command

import "context"

// Command represents a named command that can
type Command interface {
	Name() string
	Execute(ctx context.Context, args []string) error
}

type funcCommand struct {
	name     string
	function func(args []string) error
}

// FromFunc returns a command with the given name that executes the provided logic
func FromFunc(name string, function func(args []string) error) Command {
	return &funcCommand{
		name:     name,
		function: function,
	}
}

func (f *funcCommand) Name() string {
	return f.name
}

func (f *funcCommand) Execute(ctx context.Context, args []string) error {
	return f.function(args)
}
