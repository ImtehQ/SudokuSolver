package main

import (
	"bytes"
	"strings"
	"testing"
)

const testPuzzle = "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
const testSolution = "534678912672195348198342567859761423426853791713924856961537284287419635345286179"

func TestRunTextOutputUsesExactSolutionSpace(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{testPuzzle}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code=%d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{
		"Remaining valid completions: 1",
		"Unique completion: yes",
		"Next cell: r1c3",
		"4: 1 / 1 (100.000000%)",
		"Recommended digit: 4",
		"already has exactly one valid completion",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("output missing %q: %s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "randomized") {
		t.Fatalf("old randomized model leaked into output: %s", stdout.String())
	}
}

func TestRunJSONUsesExactCounts(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--json", testPuzzle}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code=%d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{`"remaining_solutions": "1"`, `"unique_solution": true`, `"recommended_digit": 4`} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("JSON missing %q: %s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "estimated_probability_percent") {
		t.Fatalf("old probability field still present: %s", stdout.String())
	}
}

func TestRunSolveUsesProbabilitySteps(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--solve", testPuzzle}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code=%d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{"Step 1: r1c3 = 4 (100.000000%", "Final grid:", testSolution[:9], testSolution[72:]} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("solve output missing %q: %s", want, stdout.String())
		}
	}
}

func TestRunRejectsInvalidPuzzle(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"123"}, strings.NewReader(""), &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "invalid puzzle") {
		t.Fatalf("run() code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}
