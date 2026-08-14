package sudoku

import "testing"

const classicPuzzle = "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
const classicSolution = "534678912672195348198342567859761423426853791713924856961537284287419635345286179"

func TestParseCompactPuzzle(t *testing.T) {
	grid, err := Parse(classicPuzzle)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if grid[0] != 5 || grid[2] != 0 || grid[80] != 9 {
		t.Fatalf("unexpected parsed cells: first=%d third=%d last=%d", grid[0], grid[2], grid[80])
	}
}

func TestParseFormattedPuzzle(t *testing.T) {
	formatted := `
5 3 . | . 7 . | . . .
6 . . | 1 9 5 | . . .
. 9 8 | . . . | . 6 .
------+-------+------
8 . . | . 6 . | . . 3
4 . . | 8 . 3 | . . 1
7 . . | . 2 . | . . 6
------+-------+------
. 6 . | . . . | 2 8 .
. . . | 4 1 9 | . . 5
. . . | . 8 . | . 7 9
`
	if _, err := Parse(formatted); err != nil {
		t.Fatalf("Parse() formatted error = %v", err)
	}
}

func TestParseRejectsDuplicate(t *testing.T) {
	bad := "550070000" + classicPuzzle[9:]
	if _, err := Parse(bad); err == nil {
		t.Fatal("Parse() accepted duplicate row value")
	}
}

func TestParseRejectsWrongCellCount(t *testing.T) {
	if _, err := Parse("123"); err == nil {
		t.Fatal("Parse() accepted short puzzle")
	}
}
