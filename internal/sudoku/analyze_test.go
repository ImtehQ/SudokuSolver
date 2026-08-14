package sudoku

import (
	"math/big"
	"testing"
)

func TestCountSolutionsClassicPuzzleIsUnique(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	count, err := CountSolutions(grid)
	if err != nil {
		t.Fatal(err)
	}
	if count.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("CountSolutions() = %s, want 1", count)
	}
}

func TestAnalyzeUniquePuzzleMakesCorrectNextDigitCertain(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	analysis, err := Analyze(grid)
	if err != nil {
		t.Fatal(err)
	}
	if analysis.RemainingSolutions != "1" || !analysis.UniqueSolution {
		t.Fatalf("unexpected solution-space result: %+v", analysis)
	}
	if analysis.NextCell == nil {
		t.Fatal("expected a next-cell analysis")
	}
	cell := analysis.NextCell
	if cell.Row != 1 || cell.Column != 3 {
		t.Fatalf("next cell = r%dc%d, want r1c3", cell.Row, cell.Column)
	}
	if cell.RecommendedDigit != 4 || cell.RecommendedProbabilityPercent != 100 || !cell.Guaranteed {
		t.Fatalf("unexpected recommendation: %+v", cell)
	}
	for _, candidate := range cell.Candidates {
		if candidate.Digit == 4 {
			if candidate.RemainingSolutions != "1" || candidate.ProbabilityPercent != 100 {
				t.Fatalf("digit 4 analysis = %+v", candidate)
			}
		} else if candidate.RemainingSolutions != "0" || candidate.ProbabilityPercent != 0 {
			t.Fatalf("non-solution candidate has non-zero weight: %+v", candidate)
		}
	}
}

func TestAnalyzeTwoSolutionPuzzleProducesFiftyFiftySplit(t *testing.T) {
	// Removing this 6/7 rectangle from a valid completed grid leaves exactly
	// two completions: the original arrangement and the swapped arrangement.
	puzzle := []byte(classicSolution)
	for _, idx := range []int{3, 4, 30, 31} {
		puzzle[idx] = '0'
	}
	grid, err := Parse(string(puzzle))
	if err != nil {
		t.Fatal(err)
	}
	analysis, err := Analyze(grid)
	if err != nil {
		t.Fatal(err)
	}
	if analysis.RemainingSolutions != "2" || analysis.UniqueSolution {
		t.Fatalf("unexpected solution count: %+v", analysis)
	}
	cell := analysis.NextCell
	if cell == nil || cell.Row != 1 || cell.Column != 4 {
		t.Fatalf("unexpected next cell: %+v", cell)
	}
	if cell.RecommendedDigit != 6 || cell.RecommendedProbabilityPercent != 50 || cell.Guaranteed {
		t.Fatalf("unexpected recommendation: %+v", cell)
	}
	for _, digit := range []uint8{6, 7} {
		candidate := findCandidate(t, *cell, digit)
		if candidate.RemainingSolutions != "1" || candidate.ProbabilityPercent != 50 {
			t.Fatalf("digit %d = %+v", digit, candidate)
		}
	}
}

func TestCountSolutionsConsistentButUnsolvablePuzzleIsZero(t *testing.T) {
	grid, err := Parse("531070000600195000098000060800060003400803001700020006060000280000419005000080079")
	if err != nil {
		t.Fatalf("puzzle should be locally consistent: %v", err)
	}
	count, err := CountSolutions(grid)
	if err != nil {
		t.Fatal(err)
	}
	if count.Sign() != 0 {
		t.Fatalf("CountSolutions() = %s, want 0", count)
	}
}

func TestSolveByProbabilitySolvesClassicPuzzleWithGuaranteedSteps(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	result, err := SolveByProbability(grid)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Solved || result.FinalGrid != classicSolution {
		t.Fatalf("unexpected solve result: solved=%v final=%q", result.Solved, result.FinalGrid)
	}
	if len(result.Steps) != 51 {
		t.Fatalf("got %d steps, want 51", len(result.Steps))
	}
	for _, step := range result.Steps {
		if !step.Cell.Guaranteed || step.Cell.RecommendedProbabilityPercent != 100 || step.RemainingSolutionsAfter != "1" {
			t.Fatalf("unique puzzle step was not guaranteed: %+v", step)
		}
	}
}

func TestSearchStateCandidateMaskMatchesGridCandidateMask(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	state := newSearchState(grid)
	for idx, value := range grid {
		if value == 0 && state.candidateMask(idx) != candidateMask(grid, idx) {
			t.Fatalf("candidate mask mismatch at cell %d", idx)
		}
	}

	idx := 2
	state.place(idx, 4)
	grid[idx] = 4
	for cell, value := range grid {
		if value == 0 && state.candidateMask(cell) != candidateMask(grid, cell) {
			t.Fatalf("candidate mask mismatch after placement at cell %d", cell)
		}
	}
}

func TestHasSolutionClassicPuzzle(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	if !HasSolution(grid) {
		t.Fatal("classic puzzle should have a solution")
	}
}

func findCandidate(t *testing.T, cell CellAnalysis, digit uint8) CandidateAnalysis {
	t.Helper()
	for _, candidate := range cell.Candidates {
		if candidate.Digit == digit {
			return candidate
		}
	}
	t.Fatalf("candidate %d not found", digit)
	return CandidateAnalysis{}
}
