package game

import "errors"

const (
	Empty = iota
	Black
	White
)

type Move struct {
	X, Y int
}

type Board struct {
	Cells [8][8]int
	Turn  int
}

func NewBoard() *Board {
	b := &Board{}
	b.Cells[3][3] = Black
	b.Cells[3][4] = White
	b.Cells[4][3] = White
	b.Cells[4][4] = Black
	b.Turn = Black
	return b
}

func (b *Board) IsValidMove(move Move) bool {
	if move.X < 0 || move.X >= 8 || move.Y < 0 || move.Y >= 8 {
		return false
	}

	if b.Cells[move.Y][move.X] != Empty {
		return false
	}

	dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}}
	for _, dir := range dirs {
		if b.isValidDirection(move, dir) {
			return true
		}
	}

	return false
}

func (b *Board) isValidDirection(move Move, dir []int) bool {
	x, y := move.X+dir[0], move.Y+dir[1]
	if x < 0 || x >= 8 || y < 0 || y >= 8 || b.Cells[y][x] != b.opponent() {
		return false
	}

	for x >= 0 && x < 8 && y >= 0 && y < 8 {
		if b.Cells[y][x] == b.Turn {
			return true
		}
		x += dir[0]
		y += dir[1]
	}

	return false
}

func (b *Board) MakeMove(move Move) error {
	if !b.IsValidMove(move) {
		return errors.New("invalid move")
	}

	dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}}
	for _, dir := range dirs {
		if b.isValidDirection(move, dir) {
			b.flipStones(move, dir)
		}
	}

	b.Cells[move.Y][move.X] = b.Turn
	b.Turn = b.opponent()

	return nil
}

func (b *Board) flipStones(move Move, dir []int) {
	x, y := move.X+dir[0], move.Y+dir[1]
	for b.Cells[y][x] == b.opponent() {
		b.Cells[y][x] = b.Turn
		x += dir[0]
		y += dir[1]
	}
}

func (b *Board) opponent() int {
	if b.Turn == Black {
		return White
	} else {
		return Black
	}
}

func (b *Board) Winner() int {
	black, white := 0, 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			switch b.Cells[y][x] {
			case Black:
				black++
			case White:
				white++
			}
		}
	}
	if black > white {
		return Black
	} else if black < white {
		return White
	} else {
		return Empty // Draw
	}
}

func (b *Board) AvailableMoves() []Move {
	var moves []Move
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			move := Move{x, y}
			if b.IsValidMove(move) {
				moves = append(moves, move)
			}
		}
	}
	return moves
}
