package sudoku

import "testing"

func TestProbabilityCounterUsesStableRowMajorMRV(t *testing.T) {
	if !newExactCounter().degreeTieBreak {
		t.Fatal("exact counter should use degree tie-breaking")
	}
	if newProbabilityCounter().degreeTieBreak {
		t.Fatal("probability counter should preserve row-major MRV for memo reuse")
	}

	grid, err := Parse("530600000002105000190302567050000423026053091013924056961000200000419000045206109")
	if err != nil {
		t.Fatal(err)
	}
	state := newSearchState(grid)
	degreeIdx, _ := selectBranchCell(&state)
	rowMajorIdx, _ := selectBranchCellRowMajor(&state)
	if degreeIdx != 6 || rowMajorIdx != 5 {
		t.Fatalf("selectors chose degree=%d row-major=%d, want 6 and 5", degreeIdx, rowMajorIdx)
	}
}
