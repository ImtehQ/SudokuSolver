package sudoku

import "testing"

func TestResetWorkspaceClearsMemoAndFramesWithoutShrinking(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	counter := newProbabilityCounter()
	if count := counter.count(grid); count.String() != "1" {
		t.Fatalf("count = %s, want 1", count)
	}
	memoSlots := len(counter.memo.entries)
	stackFrames := len(counter.stack)
	if counter.memo.used == 0 || memoSlots == 0 || stackFrames == 0 {
		t.Fatal("expected populated reusable workspace")
	}

	counter.resetWorkspace()
	if counter.memo.used != 0 {
		t.Fatalf("memo used = %d, want 0", counter.memo.used)
	}
	if len(counter.memo.entries) != memoSlots || len(counter.stack) != stackFrames {
		t.Fatal("reset should retain allocated workspace dimensions")
	}
	for idx, entry := range counter.memo.entries {
		if entry.used {
			t.Fatalf("memo entry %d still marked used", idx)
		}
	}
	for idx, frame := range counter.stack {
		if frame.initialized || frame.remainingMask != 0 || frame.start != (packedGrid{}) {
			t.Fatalf("search frame %d was not cleared", idx)
		}
	}
}
