package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/softdevstudio/oqft"
)

const defaultEpsilon = 0.01

func runStats(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("oqft stats", flag.ContinueOnError)
	flags.SetOutput(stderr)
	n := flags.Int("n", 0, "number of qubits (required)")
	cutoff := flags.Int("cutoff", defaultCutoff, "maximum distance between qubits for AQFT rotations")
	eps := flags.Float64("eps", defaultEpsilon, "add an OQFT comparison at this error bound when specified")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: oqft stats -n N [--cutoff K] [--eps E]")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "oqft stats: unexpected argument %q\n", flags.Arg(0))
		return 2
	}

	approx, err := oqft.AQFT(*n, *cutoff)
	if err != nil {
		fmt.Fprintf(stderr, "oqft stats: %v\n", err)
		return 2
	}
	exact, err := oqft.QFT(*n)
	if err != nil {
		fmt.Fprintf(stderr, "oqft stats: %v\n", err)
		return 2
	}
	report := statsReport(exact, approx, *cutoff)
	var includeOptimistic bool
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "eps" {
			includeOptimistic = true
		}
	})
	if includeOptimistic {
		circuit, err := oqft.OQFT(*n, *eps)
		if err != nil {
			fmt.Fprintf(stderr, "oqft stats: %v\n", err)
			return 2
		}
		report += fmt.Sprintf("OQFT epsilon: %g, total gates: %d, two-qubit gates: %d, depth: %d\n",
			*eps, circuit.Len(), circuit.TwoQubit(), circuit.Depth())
	}
	if _, err := io.WriteString(stdout, report); err != nil {
		fmt.Fprintf(stderr, "oqft stats: %v\n", err)
		return 1
	}
	return 0
}

func statsReport(exact, approx *oqft.Circuit, cutoff int) string {
	var report strings.Builder
	fmt.Fprintf(&report, "Qubits: %d, AQFT cutoff: %d\n", exact.N, cutoff)
	fmt.Fprintf(&report, "%-14s %10s %10s\n", "Metric", "QFT", "AQFT")
	row := func(label string, left, right int) {
		fmt.Fprintf(&report, "%-14s %10d %10d\n", label, left, right)
	}
	exactCounts, approxCounts := exact.Counts(), approx.Counts()
	for _, kind := range oqft.Kinds {
		row(kind.String(), exactCounts[kind], approxCounts[kind])
	}
	row("Total gates", exact.Len(), approx.Len())
	row("Two-qubit", exact.TwoQubit(), approx.TwoQubit())
	row("Depth", exact.Depth(), approx.Depth())
	return report.String()
}
