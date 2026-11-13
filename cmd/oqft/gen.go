package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/softdevstudio/oqft"
)

func runGen(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("oqft gen", flag.ContinueOnError)
	flags.SetOutput(stderr)
	kind := flags.String("type", "qft", "circuit type (qft or aqft)")
	n := flags.Int("n", 0, "number of qubits (required)")
	cutoff := flags.Int("cutoff", defaultCutoff, "maximum distance between qubits for AQFT rotations")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: oqft gen --type qft|aqft -n N [--cutoff K]")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "oqft gen: unexpected argument %q\n", flags.Arg(0))
		return 2
	}
	if *n < 1 {
		fmt.Fprintln(stderr, "oqft gen: -n must be a positive qubit count")
		return 2
	}

	var circuit *oqft.Circuit
	var err error
	switch *kind {
	case "qft":
		circuit, err = oqft.QFT(*n)
	case "aqft":
		circuit, err = oqft.AQFT(*n, *cutoff)
	default:
		fmt.Fprintf(stderr, "oqft gen: unknown circuit type %q\n", *kind)
		return 2
	}
	if err != nil {
		fmt.Fprintf(stderr, "oqft gen: %v\n", err)
		return 2
	}

	program, err := circuit.OpenQASM()
	if err == nil {
		_, err = io.WriteString(stdout, program)
	}
	if err != nil {
		fmt.Fprintf(stderr, "oqft gen: %v\n", err)
		return 1
	}
	return 0
}
