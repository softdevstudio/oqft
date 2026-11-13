package oqft

import (
	"math"
	"math/cmplx"
	"slices"
	"testing"
)

func TestBlockedQFTGateOrder(t *testing.T) {
	exact, err := QFT(7)
	if err != nil {
		t.Fatal(err)
	}
	for _, width := range []int{1, 2, 3, 4, 7, 8} {
		got, err := blockedQFT(exact.N, width)
		if err != nil {
			t.Fatal(err)
		}
		want := slices.DeleteFunc(slices.Clone(exact.Gates), func(g Gate) bool {
			return g.Kind == CPhase && g.Q1/width-g.Q0/width > 1
		})
		if !slices.Equal(got.Gates, want) {
			t.Errorf("width %d: gates differ from the filtered textbook circuit", width)
		}
	}
}

func TestBlockedQFTAgainstTextbook(t *testing.T) {
	const n = 7
	exact, err := QFT(n)
	if err != nil {
		t.Fatal(err)
	}
	want, err := exact.Unitary()
	if err != nil {
		t.Fatal(err)
	}
	for width := 1; width <= 4; width++ {
		circuit, err := blockedQFT(n, width)
		if err != nil {
			t.Fatal(err)
		}
		got, err := circuit.Unitary()
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
		if width == 4 {
			if distance != 0 {
				t.Errorf("two blocks should give the exact QFT, got squared error %g", distance)
			}
			continue
		}
		blocks := (n + width - 1) / width
		bound := math.Ldexp(4*math.Pi*math.Pi*float64(blocks)/3, -width)
		if !(distance > 0 && distance <= bound) {
			t.Errorf("width %d: normalized squared error %g, want between zero and %g", width, distance, bound)
		}
	}
}

func TestBlockedQFTValidation(t *testing.T) {
	for _, input := range [][2]int{{0, 2}, {7, 0}, {7, -1}} {
		if _, err := blockedQFT(input[0], input[1]); err == nil {
			t.Errorf("blockedQFT(%d, %d) accepted invalid input", input[0], input[1])
		}
	}
}
