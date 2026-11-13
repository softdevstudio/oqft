// Command blockwidths prints the block width OQFT selects for each register
// size and error target given on the command line, for the recorded results.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/softdevstudio/oqft"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: blockwidths N EPS")
		os.Exit(2)
	}
	n, err1 := strconv.Atoi(os.Args[1])
	eps, err2 := strconv.ParseFloat(os.Args[2], 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "blockwidths: bad arguments")
		os.Exit(2)
	}
	w, err := oqft.OQFTBlockWidth(n, eps)
	if err != nil {
		fmt.Fprintln(os.Stderr, "blockwidths:", err)
		os.Exit(1)
	}
	fmt.Printf("n=%d eps=%g width=%d\n", n, eps, w)
}
