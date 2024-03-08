package commands

import (
	"fmt"
	"os"
)

func CommandPrintln(a ...any) (int, error) {
	return fmt.Println(a...)
}

func CommandPrintErrorln(a ...any) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func CommandPrettyPrintln(a ...any) (int, error) {
	return fmt.Fprintln(os.Stdout, a...)
}
