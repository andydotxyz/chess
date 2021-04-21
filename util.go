package main

import (
	"log"

	"fyne.io/fyne/v2"

	"github.com/notnil/chess"
)

const preferenceKeyCurrent = "current"

func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()
	for _, m := range valid {
		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {
			return m
		}
	}

	return nil
}

func loadGameFromPreference(game *chess.Game, p fyne.Preferences) {
	cur := p.String(preferenceKeyCurrent)
	if cur == "" {
		return
	}

	load, err := chess.FEN(cur)
	if err != nil {
		log.Println("Failed to load game", err)
		return
	}
	load(game)
}

func positionToSquare(pos fyne.Position, gridSize fyne.Size) chess.Square {
	var offX, offY = -1, -1
	for x := float32(0); x <= pos.X; x += gridSize.Width / 8 {
		offX++
	}
	for y := float32(0); y <= pos.Y; y += gridSize.Height / 8 {
		offY++
	}

	return chess.Square((7-offY)*8 + offX)
}

func squareToOffset(sq chess.Square) int {
	x := sq % 8
	y := 7 - ((sq - x) / 8)

	return int(x + y*8)
}
