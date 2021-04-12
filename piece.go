package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/notnil/chess"
)

var moveStart chess.Square = chess.NoSquare

type piece struct {
	widget.Icon

	game   *chess.Game
	square chess.Square
}

func newPiece(g *chess.Game, sq chess.Square) *piece {
	p := g.Position().Board().Piece(sq)

	ret := &piece{game: g, square: sq}
	ret.ExtendBaseWidget(ret)

	ret.Resource = resourceForPiece(p)
	return ret
}

func (p *piece) Tapped(ev *fyne.PointEvent) {
	if moveStart == chess.NoSquare {
		moveStart = p.square
		return
	}

	valid := p.game.ValidMoves()
	for _, m := range valid {
		if m.S1() == moveStart && m.S2() == p.square {
			moveStart = chess.NoSquare
			move(m, p.game, grid, over)

			go func() {
				time.Sleep(time.Second)
				randomResponse(p.game)
			}()
			return
		}
	}

	message := fmt.Sprintf("Cannot move piece %d to square %v",
		p.game.Position().Board().Piece(moveStart), p.square)
	dialog.ShowInformation("Invalid move", message, win)

	moveStart = chess.NoSquare
}

func randomResponse(game *chess.Game) {
	valid := game.ValidMoves()
	m := valid[rand.Intn(len(valid))]

	move(m, game, grid, over)
}
