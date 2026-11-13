package oqft

import (
	"fmt"
	"math"
)

// OQFTWithBlockWidth constructs Algorithm 1 of arXiv:2505.00701 with an
// explicit block width. Use OQFT to select a width for an error bound.
func OQFTWithBlockWidth(n, width int) (*Circuit, error) {
	if width < 1 {
		return nil, fmt.Errorf("QFT block width must be positive, got %d", width)
	}
	c, err := NewCircuit(n)
	if err != nil {
		return nil, err
	}

	type block struct {
		low, high int
		gates     []Gate
	}
	var blocks []block
	for high := n; high > 0; {
		// Anchor the partition at qubit 0, as in blockedQFT.
		low := (high - 1) / width * width
		local, err := QFT(high - low)
		if err != nil {
			return nil, err
		}
		// Retain the local bit reversal in the wire convention.
		gates := local.Gates[:local.Len()-(high-low)/2]
		for i := range gates {
			gates[i].Q0 += low
			if gates[i].Kind.Qubits() == 2 {
				gates[i].Q1 += low
			}
		}
		blocks = append(blocks, block{low, high, gates})
		high = low
	}

	phase := func(parity int) {
		for i := parity; i+1 < len(blocks); i += 2 {
			targets, controls := blocks[i], blocks[i+1]
			for target := targets.high - 1; target >= targets.low; target-- {
				for control := controls.high - 1; control >= controls.low; control-- {
					c.CP(control, target, math.Ldexp(math.Pi, control-target))
				}
			}
		}
	}

	for i := 0; i < len(blocks); i += 2 {
		c.Gates = append(c.Gates, blocks[i].gates...)
	}
	phase(0)
	for i := 1; i < len(blocks); i++ {
		gates := blocks[i].gates
		if i%2 == 1 {
			c.Gates = append(c.Gates, gates...)
			continue
		}
		// Undo the local QFT to estimate the input in the computational basis.
		for j := len(gates) - 1; j >= 0; j-- {
			gate := gates[j]
			gate.Angle = -gate.Angle
			c.Gates = append(c.Gates, gate)
		}
	}
	phase(1)
	for i := 2; i < len(blocks); i += 2 {
		c.Gates = append(c.Gates, blocks[i].gates...)
	}
	for low := 0; low < n/2; low++ {
		c.Swap(low, n-1-low)
	}
	return c, nil
}

// OQFT returns an optimistic QFT whose normalized squared Frobenius
// error is bounded by eps. It uses exactly n qubits.
func OQFT(n int, eps float64) (*Circuit, error) {
	width, err := OQFTBlockWidth(n, eps)
	if err != nil {
		return nil, err
	}
	if width >= n/2+n%2 {
		return QFT(n)
	}
	return OQFTWithBlockWidth(n, width)
}

// OQFTBlockWidth returns the block width selected for the given error bound.
// A width that leaves at most two blocks selects the exact QFT.
func OQFTBlockWidth(n int, eps float64) (int, error) {
	if n < 1 {
		return 0, fmt.Errorf("circuit needs at least 1 qubit, got %d", n)
	}
	if !(eps > 0 && eps <= 1) {
		return 0, fmt.Errorf("epsilon must be in (0, 1], got %g", eps)
	}

	exactWidth := n/2 + n%2
	for width := 1; width < exactWidth; width++ {
		blocks := (n-1)/width + 1
		commutations := (blocks - 1) / 2
		// Appendix B bounds truncation error. For each commutation, a
		// signed cyclic phase-estimation offset of magnitude d has
		// probability at most 1/d² (including d=1). Averaged over the
		// block input, it wraps on a fraction d/2^width of inputs.
		// A wrap costs at most 4 in squared distance; otherwise the
		// cost is at most (2πd/2^width)². Summing over both signs gives
		// the bound below. The shorter leading block is never estimated.
		// Round the logarithmic bound upward to avoid subnormal products.
		truncation := 4 * math.Pi * math.Pi * float64(blocks) / 3
		commutation := 4*math.Pi*math.Pi + 8*(1+float64(width)*math.Ln2)
		// Apply the triangle inequality before squaring.
		scale := math.Sqrt(truncation) + float64(commutations)*math.Sqrt(commutation)
		logBound := math.Nextafter(2*math.Log2(scale)-float64(width), math.Inf(1))
		if logBound <= math.Log2(eps) {
			return width, nil
		}
	}
	return exactWidth, nil
}
