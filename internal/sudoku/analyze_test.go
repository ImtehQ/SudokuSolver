package sudoku

import "testing"

func TestAnalyzeIsDeterministicForSeed(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	a, err := Analyze(grid, 250, 12345, false)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Analyze(grid, 250, 12345, false)
	if err != nil {
		t.Fatal(err)
	}
	if a.SuccessfulCompletions != b.SuccessfulCompletions || a.EstimatedProbability != b.EstimatedProbability {
		t.Fatalf("same seed produced different results: %+v vs %+v", a, b)
	}
	if a.Givens != 30 || a.EmptyCells != 51 {
		t.Fatalf("unexpected grid statistics: %+v", a)
	}
}

func TestAnalyzeSolvedGridIsCertain(t *testing.T) {
	grid, err := Parse(classicSolution)
	if err != nil {
		t.Fatal(err)
	}
	analysis, err := Analyze(grid, 10, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if analysis.EstimatedProbability != 100 || analysis.SuccessfulCompletions != 10 {
		t.Fatalf("solved grid analysis = %+v", analysis)
	}
	if analysis.VerifiedSolvable == nil || !*analysis.VerifiedSolvable {
		t.Fatalf("solved grid should verify as solvable: %+v", analysis)
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

func TestHasSolutionRejectsConsistentButUnsolvablePuzzle(t *testing.T) {
	grid, err := Parse("531070000600195000098000060800060003400803001700020006060000280000419005000080079")
	if err != nil {
		t.Fatalf("puzzle should be locally consistent: %v", err)
	}
	if HasSolution(grid) {
		t.Fatal("puzzle should have no valid completion")
	}
}

func TestAnalyzeRejectsNonPositiveTrials(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Analyze(grid, 0, 1, false); err == nil {
		t.Fatal("Analyze() accepted zero trials")
	}
}
