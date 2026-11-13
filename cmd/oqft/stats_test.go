package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestStats(t *testing.T) {
	var out, diagnostics bytes.Buffer
	code := run([]string{"stats", "-n", "2", "--cutoff", "0"}, &out, &diagnostics)
	if code != 0 || diagnostics.Len() != 0 {
		t.Fatalf("exit %d, stderr %q", code, &diagnostics)
	}
	got := strings.Join(strings.Fields(out.String()), " ")
	const want = "Qubits: 2, AQFT cutoff: 0 Metric QFT AQFT h 2 2 cp 1 0 swap 1 1 rz 0 0 Total gates 4 3 Two-qubit 2 1 Depth 4 2"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStatsArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"-n", "bad"}, {"-n", "3", "extra"}, {"-n", "3", "--cutoff", "-1"}} {
		var out, diagnostics bytes.Buffer
		code := run(append([]string{"stats"}, args...), &out, &diagnostics)
		if code != 2 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", args, code, &out, &diagnostics)
		}
	}
	var out, diagnostics bytes.Buffer
	if code := run([]string{"stats", "--help"}, &out, &diagnostics); code != 0 ||
		out.Len() != 0 || !strings.Contains(diagnostics.String(), "--cutoff") {
		t.Fatalf("help: exit %d, stdout %q, stderr %q", code, &out, &diagnostics)
	}
}
