//go:generate fyne bundle -o bundled-board.go board

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/notnil/chess"
)

type boardLayout struct{}

func (b *boardLayout) Layout(cells []fyne.CanvasObject, s fyne.Size) {
	smallEdge := s.Width
	if s.Height < s.Width {
		smallEdge = s.Height
	}

	leftInset := (s.Width - smallEdge) / 2
	cellEdge := smallEdge / 8
	cellSize := fyne.NewSize(cellEdge, cellEdge)
	i := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			cells[i].Resize(cellSize)
			cells[i].Move(fyne.NewPos(
				leftInset+(float32(x)*cellEdge), float32(y)*cellEdge))

			i++
		}
	}
}

func (b *boardLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	edge := theme.IconInlineSize() * 8
	return fyne.NewSize(edge, edge)
}

func createGrid(g *chess.Game) *fyne.Container {
	var cells []fyne.CanvasObject

	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.NRGBA{0xF4, 0xE2, 0xB6, 0xFF})
			effect := canvas.NewImageFromResource(resourceOverlay1Png)
			if x%2 == y%2 {
				bg.FillColor = color.RGBA{0x73, 0x50, 0x32, 0xFF}
				effect.Resource = resourceOverlay2Png
			}

			p := newPiece(g, chess.Square(x+y*8))
			cells = append(cells, container.NewMax(bg, effect, p))
		}
	}

	return container.New(&boardLayout{}, cells...)
}

func refreshGrid(grid *fyne.Container, b *chess.Board) {
	y, x := 7, 0
	for _, cell := range grid.Objects {
		p := b.Piece(chess.Square(x + y*8))

		img := cell.(*fyne.Container).Objects[2].(*piece)
		img.Resource = resourceForPiece(p)
		img.Refresh()

		x++
		if x == 8 {
			x = 0
			y--
		}
	}
}
