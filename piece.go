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
		if m := isValidMove(p.square, chess.NoSquare, p.game); m != nil {
			moveStart = p.square
		} else {
			message := fmt.Sprintf("Cannot move piece %d",
				p.game.Position().Board().Piece(p.square))
			dialog.ShowInformation("Invalid move", message, win)
		}
		return
	}

	if m := isValidMove(moveStart, p.square, p.game); m != nil {
		moveStart = chess.NoSquare
		move(m, p.game, grid, over)

		go func() {
			time.Sleep(time.Second)
			randomResponse(p.game)
		}()
		return
	}

	message := fmt.Sprintf("Cannot move piece %d to square %v",
		p.game.Position().Board().Piece(moveStart), p.square)
	dialog.ShowInformation("Invalid move", message, win)

	moveStart = chess.NoSquare
}

func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()
	for _, m := range valid {
		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {
			return m
		}
	}

	return nil
}

func randomResponse(game *chess.Game) {
	valid := game.ValidMoves()
	if len(valid) == 0 {
		return // game ended, better control logic would avoid this
	}
	m := valid[rand.Intn(len(valid))]

	move(m, game, grid, over)
}
