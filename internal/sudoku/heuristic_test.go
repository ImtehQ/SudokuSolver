package sudoku

import (
	"math/bits"
	"testing"
)

func TestSelectBranchCellBreaksMRVTieByEmptyPeerDegree(t *testing.T) {
	grid, err := Parse("530600000002105000190302567050000423026053091013924056961000200000419000045206109")
	if err != nil {
		t.Fatal(err)
	}
	state := newSearchState(grid)
	idx, mask := selectBranchCell(&state)
	if idx != 6 {
		t.Fatalf("selectBranchCell() index = %d (r%dc%d), want 6 (r1c7)", idx, idx/Size+1, idx%Size+1)
	}
	if got := bits.OnesCount16(mask); got != 2 {
		t.Fatalf("selected candidate count = %d, want 2", got)
	}
	if state.emptyPeerDegree(6) <= state.emptyPeerDegree(5) {
		t.Fatalf("degree tie-break fixture invalid: r1c7=%d r1c6=%d", state.emptyPeerDegree(6), state.emptyPeerDegree(5))
	}
}
