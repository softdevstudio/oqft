package oqft

import (
	"math"
	"math/cmplx"
	"strconv"
	"testing"
)

func TestQFTMatchesDFT(t *testing.T) {
	for n := 1; n <= 10; n++ {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			c, err := QFT(n)
			if err != nil {
				t.Fatal(err)
			}
			counts := c.Counts()
			if counts[H] != n || counts[CPhase] != n*(n-1)/2 || counts[Swap] != n/2 {
				t.Fatalf("unexpected textbook gate counts: %v", counts)
			}
			matrix, err := c.Unitary()
			if err != nil {
				t.Fatal(err)
			}
			size := 1 << n
			if len(matrix) != size {
				t.Fatalf("got %d rows, want %d", len(matrix), size)
			}
			scale := 1 / math.Sqrt(float64(size))
			for output, row := range matrix {
				if len(row) != size {
					t.Fatalf("row %d has length %d, want %d", output, len(row), size)
				}
				for input, got := range row {
					angle := 2 * math.Pi * float64(output*input%size) / float64(size)
					want := cmplx.Rect(scale, angle)
					if diff := cmplx.Abs(got - want); !(diff < 1e-12) {
						t.Fatalf("matrix[%d][%d] = %v, want %v (error %g)", output, input, got, want, diff)
					}
				}
			}
		})
	}
}

func TestUnitaryBasisOrder(t *testing.T) {
	c := Circuit{N: 1}
	c.H(0).RZ(0, math.Pi)
	matrix, err := c.Unitary()
	if err != nil {
		t.Fatal(err)
	}
	const amplitude = math.Sqrt2 / 2
	checkAmplitudes(t, matrix[0], []complex128{-amplitude * 1i, -amplitude * 1i})
	checkAmplitudes(t, matrix[1], []complex128{amplitude * 1i, -amplitude * 1i})
}

func TestQFTAndUnitaryValidation(t *testing.T) {
	if _, err := QFT(0); err == nil {
		t.Error("QFT accepted zero qubits")
	}
	for _, c := range []Circuit{
		{N: 0},
		{N: MaxUnitaryQubits + 1},
		{N: 1, Gates: []Gate{{Kind: H, Q0: 1}}},
	} {
		if _, err := c.Unitary(); err == nil {
			t.Errorf("Unitary accepted %+v", c)
		}
	}
}
