package sudoku

import (
	"fmt"
	"strings"
	"unicode"
)

const Size = 9
const CellCount = Size * Size

// Grid stores a Sudoku grid row-major. Zero represents an empty cell.
type Grid [CellCount]uint8

// Parse accepts 1-9 as givens, 0 or . as blanks, and ignores common visual separators.
func Parse(input string) (Grid, error) {
	var grid Grid
	cells := make([]uint8, 0, CellCount)

	for _, r := range input {
		switch {
		case r >= '1' && r <= '9':
			cells = append(cells, uint8(r-'0'))
		case r == '0' || r == '.':
			cells = append(cells, 0)
		case unicode.IsSpace(r) || strings.ContainsRune("|+-_,", r):
			// Allow human-friendly formatting.
		default:
			return grid, fmt.Errorf("unexpected character %q", r)
		}
	}

	if len(cells) != CellCount {
		return grid, fmt.Errorf("expected 81 cells, got %d", len(cells))
	}
	copy(grid[:], cells)

	if err := Validate(grid); err != nil {
		return grid, err
	}
	return grid, nil
}

// Validate checks the currently filled cells for direct row, column, and box conflicts.
func Validate(grid Grid) error {
	for i, value := range grid {
		if value > 9 {
			return fmt.Errorf("cell %d has invalid value %d", i, value)
		}
	}

	for row := 0; row < Size; row++ {
		if err := validateUnit(grid, rowCells(row), fmt.Sprintf("row %d", row+1)); err != nil {
			return err
		}
	}
	for col := 0; col < Size; col++ {
		if err := validateUnit(grid, colCells(col), fmt.Sprintf("column %d", col+1)); err != nil {
			return err
		}
	}
	for box := 0; box < Size; box++ {
		if err := validateUnit(grid, boxCells(box), fmt.Sprintf("box %d", box+1)); err != nil {
			return err
		}
	}
	return nil
}

// Compact returns the grid as 81 digits with 0 representing empty cells.
func Compact(grid Grid) string {
	var b strings.Builder
	b.Grow(CellCount)
	for _, value := range grid {
		b.WriteByte('0' + byte(value))
	}
	return b.String()
}

func validateUnit(grid Grid, cells [Size]int, label string) error {
	var seen [10]bool
	for _, idx := range cells {
		value := grid[idx]
		if value == 0 {
			continue
		}
		if seen[value] {
			return fmt.Errorf("duplicate %d in %s", value, label)
		}
		seen[value] = true
	}
	return nil
}

func rowCells(row int) [Size]int {
	var cells [Size]int
	for col := 0; col < Size; col++ {
		cells[col] = row*Size + col
	}
	return cells
}

func colCells(col int) [Size]int {
	var cells [Size]int
	for row := 0; row < Size; row++ {
		cells[row] = row*Size + col
	}
	return cells
}

func boxCells(box int) [Size]int {
	var cells [Size]int
	startRow := (box / 3) * 3
	startCol := (box % 3) * 3
	k := 0
	for row := startRow; row < startRow+3; row++ {
		for col := startCol; col < startCol+3; col++ {
			cells[k] = row*Size + col
			k++
		}
	}
	return cells
}
