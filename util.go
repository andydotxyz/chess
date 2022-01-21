package main

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"time"

	"fyne.io/fyne/v2"

	"github.com/notnil/chess"
)

const PREFERENCE_KEY_CURRENT = "current"

func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()
	for _, m := range valid {
		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {
			return m
		}
	}

	return nil
}

func positionToSquare(pos fyne.Position, gridSize fyne.Size) chess.Square {
	var offX, offY = -1, -1
	cellEdge := cellSize(gridSize)
	for x := (gridSize.Width - cellEdge*8) / 2; x <= pos.X; x += cellEdge {
		offX++
	}
	for y := float32(0); y <= pos.Y; y += cellEdge {
		offY++
	}

	return chess.Square((7-offY)*8 + offX)
}

func squareToOffset(sq chess.Square) int {
	x := sq % 8
	y := 7 - ((sq - x) / 8)

	return int(x + y*8)
}

func move1(m *chess.Move, game *chess.Game, u *ui, notHuman bool) {
	if m == nil {
		return
	}

	off := squareToOffset(m.S1())
	cell := u.grid.objects[off].(*fyne.Container)
	u.over.Move(cell.Position())

	img := cell.Objects[2].(*piece)

	u.over.Resource = resourceForPiece(game.Position().Board().Piece(m.S1()))
	u.over.Resize(img.Size())
	u.over.Refresh() // clear our old resource before showing

	u.over.Show()
	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	cell = u.grid.objects[off].(*fyne.Container)
	pos2 := cell.Position()

	if notHuman {
		a := canvas.NewPositionAnimation(u.over.Position(), pos2, time.Millisecond*500, func(p fyne.Position) {
			u.over.Move(p)
			u.over.Refresh()
		})
		a.Start()
		time.Sleep(500 * time.Millisecond)
	}
	u.over.Hide()

}

func move2(m *chess.Move, game *chess.Game, u *ui) {
	white := game.Position().Turn() == chess.White
	_ = u.blackTurn.Set(white)
	_ = u.outcome.Set(string(game.Outcome()))

	if game.Outcome() != chess.NoOutcome {
		result := "draw"
		switch game.Outcome().String() {
		case "1-0":
			result = "won"
		case "0-1":
			result = "lost"
		}

		dialog.ShowInformation("Game ended",
			"Game "+result+" because "+game.Method().String(), u.win)
	}
	u.refreshGrid(game)
}
