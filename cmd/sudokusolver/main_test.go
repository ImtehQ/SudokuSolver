package main

import (
	"bytes"
	"strings"
	"testing"
)

const testPuzzle = "530070000600195000098000060800060003400803001700020006060000280000419005000080079"

func TestRunTextOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--trials", "20", "--seed", "7", "--verify", testPuzzle}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code=%d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{"Estimated solution probability:", "Exact solvability: yes", "empirical randomized completion score"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("output missing %q: %s", want, stdout.String())
		}
	}
}

func TestRunJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--json", "--trials", "10", "--seed", "1", testPuzzle}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"estimated_probability_percent"`) {
		t.Fatalf("unexpected JSON output: %s", stdout.String())
	}
}

func TestRunRejectsInvalidPuzzle(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"123"}, strings.NewReader(""), &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "invalid puzzle") {
		t.Fatalf("run() code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}
