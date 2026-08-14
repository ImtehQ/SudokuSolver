package sudoku

import "math/bits"

var cellPeers = buildCellPeers()

func buildCellPeers() [Size * Size][20]uint8 {
	var peers [Size * Size][20]uint8
	for idx := 0; idx < Size*Size; idx++ {
		row := idx / Size
		col := idx % Size
		box := (row/3)*3 + col/3
		next := 0
		for peer := 0; peer < Size*Size; peer++ {
			if peer == idx {
				continue
			}
			peerRow := peer / Size
			peerCol := peer % Size
			peerBox := (peerRow/3)*3 + peerCol/3
			if peerRow == row || peerCol == col || peerBox == box {
				peers[idx][next] = uint8(peer)
				next++
			}
		}
	}
	return peers
}

func (s *searchState) emptyPeerDegree(idx int) int {
	degree := 0
	for _, peer := range cellPeers[idx] {
		if s.grid[peer] == 0 {
			degree++
		}
	}
	return degree
}

func selectBranchCell(s *searchState) (int, uint16) {
	bestIdx := -1
	bestMask := uint16(0)
	bestCount := 10
	bestDegree := -1
	for idx, digit := range s.grid {
		if digit != 0 {
			continue
		}
		mask := s.candidateMask(idx)
		count := bits.OnesCount16(mask)
		if count > bestCount {
			continue
		}
		degree := s.emptyPeerDegree(idx)
		if count < bestCount || degree > bestDegree {
			bestIdx = idx
			bestMask = mask
			bestCount = count
			bestDegree = degree
		}
	}
	return bestIdx, bestMask
}
