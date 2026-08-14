package main

import (
	"strings"
	"testing"
)

func TestNormalizePuzzleLine(t *testing.T) {
	puzzle := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	normalized, ok, err := normalizePuzzleLine(puzzle + " extra-metadata")
	if err != nil {
		t.Fatalf("normalizePuzzleLine returned error: %v", err)
	}
	if !ok {
		t.Fatal("normalizePuzzleLine skipped a valid puzzle")
	}
	if len(normalized) != 81 {
		t.Fatalf("normalized length = %d, want 81", len(normalized))
	}
	if strings.Contains(normalized, ".") {
		t.Fatalf("normalized puzzle still contains dots: %q", normalized)
	}
}

func TestReadPuzzlesSkipsCommentsAndHonorsLimit(t *testing.T) {
	puzzle := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
	input := "# source\n\n" + puzzle + "\n" + puzzle + "\n"
	puzzles, err := readPuzzles(strings.NewReader(input), 1)
	if err != nil {
		t.Fatalf("readPuzzles returned error: %v", err)
	}
	if len(puzzles) != 1 {
		t.Fatalf("got %d puzzles, want 1", len(puzzles))
	}
}

func TestPercentileInterpolates(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if got := percentile(values, 0.50); got != 3 {
		t.Fatalf("median = %v, want 3", got)
	}
	if got := percentile(values, 0.95); got != 4.8 {
		t.Fatalf("p95 = %v, want 4.8", got)
	}
}
