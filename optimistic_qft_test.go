package oqft

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestOptimisticQFTCommutation(t *testing.T) {
	for _, tc := range []struct {
		n, width int
		exact    bool
	}{
		{1, 3, true}, {5, 3, true}, {6, 2, true}, {7, 2, false}, {9, 2, false},
	} {
		baseline, err := blockedQFT(tc.n, tc.width)
		if err != nil {
			t.Fatal(err)
		}
		optimistic, err := optimisticQFT(tc.n, tc.width)
		if err != nil {
			t.Fatal(err)
		}
		if optimistic.N != tc.n {
			t.Fatalf("got %d qubits, want %d", optimistic.N, tc.n)
		}
		want, err := baseline.Unitary()
		if err != nil {
			t.Fatal(err)
		}
		got, err := optimistic.Unitary()
		if err != nil {
			t.Fatal(err)
		}
		var squaredError float64
		for output, row := range got {
			for input, amplitude := range row {
				delta := cmplx.Abs(amplitude - want[output][input])
				squaredError += delta * delta
			}
		}
		distance := squaredError / float64(len(want))
		if tc.exact && !(distance < 1e-24) {
			t.Errorf("%+v: normalized squared error %g, want zero", tc, distance)
		}
		if !tc.exact && !(distance > 1e-6 && distance < 0.5) {
			t.Errorf("%+v: normalized squared error %g, want a small nonzero error", tc, distance)
		}
	}
}

func TestOptimisticQFTDepth(t *testing.T) {
	const width = 4
	small, err := optimisticQFT(8*width, width)
	if err != nil {
		t.Fatal(err)
	}
	large, err := optimisticQFT(512*width, width)
	if err != nil {
		t.Fatal(err)
	}
	if large.Depth() != small.Depth() || large.Depth() > 11*width {
		t.Fatalf("depth grew with the number of blocks or exceeded the stage bound: %d, %d", small.Depth(), large.Depth())
	}
}

func TestOptimisticQFTValidation(t *testing.T) {
	for _, input := range [][2]int{{0, 2}, {7, 0}, {7, -1}} {
		if _, err := optimisticQFT(input[0], input[1]); err == nil {
			t.Errorf("optimisticQFT(%d, %d) accepted invalid input", input[0], input[1])
		}
	}
}

func TestOptimisticQFTThreeSingleQubitBlocks(t *testing.T) {
	c, err := optimisticQFT(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	// The middle phase couples blocks 0 and 1. The last phase uses
	// block 2's inverse QFT to control the output on block 1.
	want := []Gate{
		{Kind: H, Q0: 2}, {Kind: H, Q0: 0},
		{Kind: CPhase, Q0: 1, Q1: 2, Angle: math.Pi / 2},
		{Kind: H, Q0: 1}, {Kind: H, Q0: 0},
		{Kind: CPhase, Q0: 0, Q1: 1, Angle: math.Pi / 2},
		{Kind: H, Q0: 0}, {Kind: Swap, Q0: 0, Q1: 2},
	}
	if len(c.Gates) != len(want) {
		t.Fatalf("got %d gates, want %d", len(c.Gates), len(want))
	}
	for i, gate := range c.Gates {
		if gate != want[i] {
			t.Errorf("gate %d: got %+v, want %+v", i, gate, want[i])
		}
	}
}
