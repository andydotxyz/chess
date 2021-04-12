package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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
