package commands

import (
	"fmt"
	"os"
)

func CommandPrintln(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func CommandPrintErrorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func CommandPrettyPrintln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stdout, a...)
}
