// Command oqft is the OQFT command line tool.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/softdevstudio/oqft"
)

const usage = `oqft: optimistic quantum Fourier transform toolkit

Usage:
  oqft version    print the version
  oqft help       show this help

More commands arrive phase by phase; see CLAUDE.md.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "oqft %s\n", oqft.Version)
		return 0
	case "help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "oqft: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
