package benchmarks_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ImtehQ/SudokuSolver/internal/sudoku"
)

type fixture struct {
	Difficulty string `json:"difficulty"`
	Puzzle     string `json:"puzzle"`
}

func TestBenchmarkPuzzlesAreUniqueAndSolvable(t *testing.T) {
	data, err := os.ReadFile("puzzles.json")
	if err != nil {
		t.Fatal(err)
	}

	var fixtures []fixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures) != 4 {
		t.Fatalf("expected 4 benchmark fixtures, got %d", len(fixtures))
	}

	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(fixture.Difficulty, func(t *testing.T) {
			grid, err := sudoku.Parse(fixture.Puzzle)
			if err != nil {
				t.Fatalf("invalid benchmark puzzle: %v", err)
			}

			count, err := sudoku.CountSolutions(grid)
			if err != nil {
				t.Fatal(err)
			}
			if count.String() != "1" {
				t.Fatalf("expected exactly one solution, got %s", count.String())
			}

			result, err := sudoku.SolveByProbability(grid)
			if err != nil {
				t.Fatal(err)
			}
			if !result.Solved {
				t.Fatal("probability-guided solver did not finish the benchmark puzzle")
			}
			for _, step := range result.Steps {
				if !step.Cell.Guaranteed {
					t.Fatalf("step %d was not guaranteed on a unique puzzle", step.Step)
				}
			}
		})
	}
}
