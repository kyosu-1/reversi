package strategy

import (
	"math/rand"
	"time"

	"github.com/kyosu-1/reversi/server/internal/game"
)

type Strategy interface {
	BestMove(board *game.Board) game.Move
}

type RandomStrategy struct{}

func (rs *RandomStrategy) BestMove(board *game.Board) game.Move {
	rand.Seed(time.Now().UnixNano())
	moves := board.AvailableMoves()
	return moves[rand.Intn(len(moves))]
}
