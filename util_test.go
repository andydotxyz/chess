package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestIsValidMove(t *testing.T) {
	g := chess.NewGame()
	m := isValidMove(chess.A2, chess.A3, g)
	assert.NotNil(t, m)
	assert.Equal(t, chess.A2, m.S1())
	assert.Equal(t, chess.A3, m.S2())
	m = isValidMove(chess.A2, chess.A4, g)
	assert.NotNil(t, m)
	assert.Equal(t, chess.A2, m.S1())
	assert.Equal(t, chess.A4, m.S2())
	assert.Nil(t, isValidMove(chess.A2, chess.A5, g))
}

func TestPositionToSquare(t *testing.T) {
	assert.Equal(t, chess.A8, positionToSquare(fyne.NewPos(5, 5), fyne.NewSize(80, 80)))
	assert.Equal(t, chess.H8, positionToSquare(fyne.NewPos(75, 5), fyne.NewSize(80, 80)))
	assert.Equal(t, chess.A1, positionToSquare(fyne.NewPos(5, 75), fyne.NewSize(80, 80)))
}

func TestSquareToOffset(t *testing.T) {
	assert.Equal(t, 0, squareToOffset(chess.A8))
	assert.Equal(t, 7, squareToOffset(chess.H8))
	assert.Equal(t, 56, squareToOffset(chess.A1))
}
