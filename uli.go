// Package main contains the main CLI
package main

import (
	"context"
	"os"

	"github.com/dimes/uli/sports/nhl"
	"github.com/dimes/uli/util/command"
)

func main() {
	commands := []command.Command{
		nhl.NewNHL(),
	}
	set, err := command.NewSet(commands)

	if err != nil {
		panic(err)
	}

	if err := set.Execute(context.Background(), os.Args[1:]); err != nil {
		panic(err)
	}
}
