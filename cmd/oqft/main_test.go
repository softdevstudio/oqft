package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/softdevstudio/oqft"
)

func TestVersion(t *testing.T) {
	var out, errb bytes.Buffer
	if code := run([]string{"version"}, &out, &errb); code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errb.String())
	}
	if got, want := out.String(), "oqft "+oqft.Version+"\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestUsageErrors(t *testing.T) {
	for _, args := range [][]string{nil, {"nope"}} {
		var out, errb bytes.Buffer
		if code := run(args, &out, &errb); code != 2 {
			t.Errorf("%v: exit %d, want 2", args, code)
		}
		if !strings.Contains(errb.String(), "Usage:") {
			t.Errorf("%v: no usage on stderr", args)
		}
	}
}
