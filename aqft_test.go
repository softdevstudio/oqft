package oqft

import (
	"math"
	"math/cmplx"
	"slices"
	"testing"
)

func TestAQFTGateCutoff(t *testing.T) {
	const n = 6
	exact, err := QFT(n)
	if err != nil {
		t.Fatal(err)
	}
	for cutoff := 0; cutoff <= n; cutoff++ {
		c, err := AQFT(n, cutoff)
		if err != nil {
			t.Fatal(err)
		}
		want := slices.DeleteFunc(slices.Clone(exact.Gates), func(g Gate) bool {
			return g.Kind == CPhase && g.Q1-g.Q0 > cutoff
		})
		if !slices.Equal(c.Gates, want) {
			t.Fatalf("cutoff %d produced %v, want %v", cutoff, c.Gates, want)
		}
	}
}

func TestAQFTErrorBoundDecreases(t *testing.T) {
	for n := 2; n <= 10; n++ {
		exact, err := QFT(n)
		if err != nil {
			t.Fatal(err)
		}
		want, err := exact.Unitary()
		if err != nil {
			t.Fatal(err)
		}
		previousDistance, previousBound := math.Inf(1), math.Inf(1)
		for cutoff := 0; cutoff < n; cutoff++ {
			c, err := AQFT(n, cutoff)
			if err != nil {
				t.Fatal(err)
			}
			got, err := c.Unitary()
			if err != nil {
				t.Fatal(err)
			}
			var squaredError float64
			for row := range want {
				for col, amplitude := range want[row] {
					delta := cmplx.Abs(got[row][col] - amplitude)
					squaredError += delta * delta
				}
			}
			// Normalize by the squared Frobenius norm of a unitary.
			distance := math.Sqrt(squaredError / float64(len(want)))
			// Every omitted rotation contributes at most its angle to
			// the phase difference of any matrix entry.
			var bound float64
			for separation := cutoff + 1; separation < n; separation++ {
				bound += float64(n-separation) * math.Ldexp(math.Pi, -separation)
			}
			if !(distance <= bound+1e-12 && bound < previousBound) {
				t.Fatalf("%d qubits, cutoff %d: distance %g, bound %g, previous bound %g", n, cutoff, distance, bound, previousBound)
			}
			// Chord distance increases with phase difference up to π.
			// Beyond that, phase wrapping can reverse the comparison.
			if previousBound <= math.Pi && !(distance < previousDistance) {
				t.Fatalf("%d qubits, cutoff %d: distance %g did not decrease from %g", n, cutoff, distance, previousDistance)
			}
			previousDistance, previousBound = distance, bound
		}
		if previousDistance != 0 {
			t.Fatalf("%d qubits: full cutoff error %g, want zero", n, previousDistance)
		}
	}
}

func TestAQFTValidation(t *testing.T) {
	for _, args := range [][2]int{{0, 0}, {2, -1}} {
		if _, err := AQFT(args[0], args[1]); err == nil {
			t.Errorf("AQFT(%d, %d) accepted invalid arguments", args[0], args[1])
		}
	}
}
