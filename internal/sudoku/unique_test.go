package sudoku

import "testing"

func TestFindOneSolutionReturnsClassicCompletion(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	solution, ok := findOneSolution(grid)
	if !ok {
		t.Fatal("expected to find classic puzzle solution")
	}
	if got := Compact(solution); got != classicSolution {
		t.Fatalf("solution = %q, want %q", got, classicSolution)
	}
}

func TestUniqueProbabilityTraceMatchesExactCandidateSemantics(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	result, err := SolveByProbability(grid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Steps) == 0 {
		t.Fatal("expected probability steps")
	}

	first := result.Steps[0]
	if first.Cell.Row != 1 || first.Cell.Column != 3 || first.Cell.RecommendedDigit != 4 {
		t.Fatalf("unexpected first step: %+v", first)
	}
	if first.RemainingSolutionsBefore != "1" || first.RemainingSolutionsAfter != "1" || !first.Cell.Guaranteed {
		t.Fatalf("unexpected unique-step solution counts: %+v", first)
	}
	for _, candidate := range first.Cell.Candidates {
		if candidate.Digit == 4 {
			if !candidate.LocallyAllowed || candidate.RemainingSolutions != "1" || candidate.ProbabilityPercent != 100 {
				t.Fatalf("solution candidate = %+v", candidate)
			}
		} else if candidate.RemainingSolutions != "0" || candidate.ProbabilityPercent != 0 {
			t.Fatalf("non-solution candidate = %+v", candidate)
		}
	}
}

func TestProbabilitySolverKeepsNonUniqueExactPath(t *testing.T) {
	puzzle := []byte(classicSolution)
	for _, idx := range []int{3, 4, 30, 31} {
		puzzle[idx] = '0'
	}
	grid, err := Parse(string(puzzle))
	if err != nil {
		t.Fatal(err)
	}
	result, err := SolveByProbability(grid)
	if err != nil {
		t.Fatal(err)
	}
	if result.InitialSolutions != "2" {
		t.Fatalf("initial solutions = %s, want 2", result.InitialSolutions)
	}
	if len(result.Steps) == 0 {
		t.Fatal("expected solve steps")
	}
	first := result.Steps[0]
	if first.RemainingSolutionsBefore != "2" || first.RemainingSolutionsAfter != "1" {
		t.Fatalf("unexpected first branch counts: %+v", first)
	}
	if first.Cell.Guaranteed || first.Cell.RecommendedProbabilityPercent != 50 {
		t.Fatalf("two-solution first step should be a 50%% non-guaranteed choice: %+v", first.Cell)
	}
}
