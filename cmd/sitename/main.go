package main

import (
	"os"

	"github.com/sitename/sitename/cmd/sitename/commands"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
