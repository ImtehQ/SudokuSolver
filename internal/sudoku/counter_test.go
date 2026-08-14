package sudoku

import (
	"math"
	"math/big"
	"testing"
)

func TestExactCountPromotesPastUint64(t *testing.T) {
	count := exactCount{small: math.MaxUint64}
	count.add(oneExactCount())

	want := new(big.Int).Lsh(big.NewInt(1), 64)
	if got := count.bigInt(); got.Cmp(want) != 0 {
		t.Fatalf("promoted count = %s, want %s", got, want)
	}
	if count.large == nil {
		t.Fatal("expected uint64 overflow to promote to big.Int")
	}
}

func TestMemoTableUsesFixedWindow(t *testing.T) {
	memo := newMemoTable(64)
	for value := 1; value <= 100; value++ {
		key := packedGrid{uint64(value)}
		memo.set(key, exactCount{small: uint64(value)})
	}

	if len(memo.entries) != 64 {
		t.Fatalf("memo slots = %d, want 64", len(memo.entries))
	}
	if memo.used != memo.limit {
		t.Fatalf("memo used = %d, want bounded limit %d", memo.used, memo.limit)
	}
	if memo.limit != 48 {
		t.Fatalf("memo load limit = %d, want 48", memo.limit)
	}
}

func TestExactCounterPreallocatesOneFramePerPossibleDepth(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}

	empty := 0
	for _, digit := range grid {
		if digit == 0 {
			empty++
		}
	}
	counter := newExactCounter()
	if count := counter.count(grid); count.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("count = %s, want 1", count)
	}
	if len(counter.stack) != empty+1 {
		t.Fatalf("workspace frames = %d, want %d", len(counter.stack), empty+1)
	}
	if slots := len(counter.memo.entries); slots != exactMemoSmallSlots && slots != exactMemoLargeSlots {
		t.Fatalf("memo window = %d slots, want %d or %d", slots, exactMemoSmallSlots, exactMemoLargeSlots)
	}
}

func TestSearchStatePackedKeyTracksPlacements(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatal(err)
	}
	state := newSearchState(grid)
	state.place(2, 4)
	grid[2] = 4
	if state.key != packGrid(grid) {
		t.Fatal("packed memo key did not track placement")
	}
}
