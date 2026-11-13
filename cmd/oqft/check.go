package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/softdevstudio/oqft"
)

func runCheck(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("oqft check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	n := flags.Int("n", 0, "number of qubits (required)")
	eps := flags.Float64("eps", defaultEpsilon, "OQFT mean squared error bound")
	explicitWidth := flags.Int("block-width", 0, "diagnostic block width; zero selects by epsilon")
	threshold := flags.Float64("threshold", 0.1, "flag states with squared error above this value")
	basis := -1
	flags.Func("basis", "check only this decimal basis index", func(value string) error {
		index, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if index < 0 {
			return fmt.Errorf("basis index must be nonnegative")
		}
		basis = index
		return nil
	})
	samples := flags.Int("samples", 0, "number of distinct random basis states")
	seed := flags.Uint64("seed", 1, "seed for random sampling")
	format := flags.String("format", "text", "report format (text or csv)")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: oqft check -n N [--eps E] [--basis X | --samples K] [--format text|csv]")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	fail := func(message string) int {
		fmt.Fprintln(stderr, "oqft check:", message)
		return 2
	}
	switch {
	case flags.NArg() != 0:
		return fail("unexpected positional argument")
	case *n < 1 || *n > oqft.MaxSimulationQubits:
		return fail(fmt.Sprintf("-n must be in 1..%d", oqft.MaxSimulationQubits))
	case *explicitWidth < 0 || *explicitWidth > *n:
		return fail("--block-width must be between zero and the register size")
	case !(*threshold >= 0 && *threshold <= 4):
		return fail("--threshold must be in [0, 4]")
	case *format != "text" && *format != "csv":
		return fail("--format must be text or csv")
	}
	size := 1 << *n
	switch {
	case basis >= size:
		return fail("--basis is outside the register")
	case *samples < 0 || *samples > size:
		return fail("--samples must be between zero and the number of basis states")
	case basis >= 0 && *samples != 0:
		return fail("choose either --basis or --samples")
	case basis == -1 && *samples == 0 && *n > oqft.MaxUnitaryQubits:
		return fail("use --basis or --samples for registers above 10 qubits")
	}

	width := *explicitWidth
	var circuit *oqft.Circuit
	var err error
	if width == 0 {
		width, err = oqft.OQFTBlockWidth(*n, *eps)
		if err != nil {
			return fail(err.Error())
		}
		circuit, err = oqft.OQFT(*n, *eps)
	} else {
		circuit, err = oqft.OQFTWithBlockWidth(*n, width)
	}
	if err != nil {
		return fail(err.Error())
	}
	checker, err := oqft.NewChecker(circuit)
	if err != nil {
		return fail(err.Error())
	}

	count := size
	var selected []int
	switch {
	case basis >= 0:
		count, selected = 1, []int{basis}
	case *samples > 0 && *samples < size:
		count = *samples
		// Draw without replacement with a sparse partial Fisher-Yates shuffle.
		rng := rand.New(rand.NewPCG(*seed, 0))
		swaps := make(map[int]int, count)
		selected = make([]int, count)
		for i := range selected {
			j := i + rng.IntN(size-i)
			value, ok := swaps[j]
			if !ok {
				value = j
			}
			selected[i] = value
			value, ok = swaps[i]
			if !ok {
				value = i
			}
			swaps[j] = value
			delete(swaps, i)
		}
	}

	if err := writeCheck(stdout, checker, *n, width, count, selected, *threshold, *format); err != nil {
		fmt.Fprintln(stderr, "oqft check:", err)
		return 1
	}
	return 0
}

func writeCheck(out io.Writer, checker *oqft.Checker, n, width, count int, selected []int, threshold float64, format string) error {
	csvOut := csv.NewWriter(out)
	header := []string{"basis", "bits", "estimated_blocks", "squared_error", "failed"}
	write := func(fields []string) error {
		if format == "csv" {
			return csvOut.Write(fields)
		}
		_, err := fmt.Fprintln(out, strings.Join(fields, "\t"))
		return err
	}
	if err := write(header); err != nil {
		return err
	}
	var sum float64
	var failures int
	for i := 0; i < count; i++ {
		input := i
		if selected != nil {
			input = selected[i]
		}
		distance, err := checker.StateError(input)
		if err != nil {
			return err
		}
		sum += distance
		failed := distance > threshold
		if failed {
			failures++
		}
		fields := []string{strconv.Itoa(input), fmt.Sprintf("%0*b", n, input),
			estimatedBlocks(n, width, input), strconv.FormatFloat(distance, 'g', -1, 64), strconv.FormatBool(failed)}
		if err := write(fields); err != nil {
			return err
		}
	}
	if format == "csv" {
		csvOut.Flush()
		return csvOut.Error()
	}
	_, err := fmt.Fprintf(out, "Checked %d of %d states, failures %d (threshold %g), mean squared error %g\n",
		count, 1<<n, failures, threshold, sum/float64(count))
	return err
}

func estimatedBlocks(n, width, input int) string {
	blocks := (n-1)/width + 1
	var values []string
	// Block indices run from high to low. The last block has no
	// outgoing phase, so its inverse and reapplication cancel.
	for block := 2; block+1 < blocks; block += 2 {
		low := (blocks - 1 - block) * width
		value := (input >> low) & ((1 << width) - 1)
		values = append(values, fmt.Sprintf("%d=%0*b", block, width, value))
	}
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, " ")
}
