package main

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestCheckCLI(t *testing.T) {
	invoke := func(args ...string) string {
		t.Helper()
		var out, errOut bytes.Buffer
		if code := run(append([]string{"check"}, args...), &out, &errOut); code != 0 || errOut.Len() != 0 {
			t.Fatalf("exit %d, stderr %q", code, &errOut)
		}
		return out.String()
	}
	text := invoke("-n", "2")
	if !strings.Contains(text, "Checked 4 of 4 states, failures 0") {
		t.Fatalf("unexpected exhaustive report: %s", text)
	}
	args := []string{"-n", "3", "--samples", "5", "--seed", "42", "--format", "csv"}
	first := invoke(args...)
	if first != invoke(args...) {
		t.Fatal("sample is not reproducible")
	}
	rows, err := csv.NewReader(strings.NewReader(first)).ReadAll()
	if err != nil || len(rows) != 6 {
		t.Fatalf("CSV rows %v, error %v", rows, err)
	}
	seen := make(map[string]bool)
	for _, row := range rows[1:] {
		if seen[row[0]] || row[3] != "0" || row[4] != "false" {
			t.Fatalf("unexpected sample row %v", row)
		}
		seen[row[0]] = true
	}
	if got := invoke("-n", "3", "--basis", "7", "--format", "csv"); !strings.Contains(got, "\n7,111,-,0,false\n") {
		t.Fatalf("unexpected explicit basis report: %s", got)
	}
	if got := invoke("-n", "4", "--basis", "010", "--format", "csv"); !strings.Contains(got, "\n10,1010,-,0,false\n") {
		t.Fatalf("unexpected decimal basis report: %s", got)
	}
}

func TestCheckWraparoundMap(t *testing.T) {
	for _, tc := range []struct{ basis, bits, block, failed string }{
		{"55", "000000110111", "2=110", "false"},
		{"63", "000000111111", "2=111", "true"},
	} {
		var out, diagnostics bytes.Buffer
		args := []string{"check", "-n", "12", "--block-width", "3", "--basis", tc.basis, "--threshold", "0.5", "--format", "csv"}
		if code := run(args, &out, &diagnostics); code != 0 {
			t.Fatalf("exit %d: %s", code, &diagnostics)
		}
		rows, err := csv.NewReader(&out).ReadAll()
		if err != nil || len(rows) != 2 {
			t.Fatalf("CSV rows %v, error %v", rows, err)
		}
		row := rows[1]
		if row[0] != tc.basis || row[1] != tc.bits || row[2] != tc.block || row[4] != tc.failed {
			t.Errorf("unexpected wraparound row: %v", row)
		}
	}
}

func TestCheckInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		nil, {"-n", "21"}, {"-n", "11"}, {"-n", "2", "--basis", "4"},
		{"-n", "2", "--samples", "5"}, {"-n", "2", "--basis", "1", "--samples", "1"},
		{"-n", "2", "--format", "json"}, {"-n", "2", "--threshold", "NaN"},
		{"-n", "2", "--eps", "0"},
	} {
		var out, errOut bytes.Buffer
		if code := run(append([]string{"check"}, args...), &out, &errOut); code != 2 || out.Len() != 0 || errOut.Len() == 0 {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", args, code, &out, &errOut)
		}
	}
}
