// Command oqft is the OQFT command line tool.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/softdevstudio/oqft"
)

const defaultCutoff = 2

const usage = `oqft: optimistic quantum Fourier transform toolkit

Usage:
  oqft version    print the version
  oqft gen        export a circuit as OpenQASM 3
  oqft stats      compare QFT and AQFT gate counts and depth
  oqft check      measure per-state error against the textbook QFT
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
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "stats":
		return runStats(args[1:], stdout, stderr)
	case "gen":
		return runGen(args[1:], stdout, stderr)
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
