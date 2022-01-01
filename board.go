//go:generate fyne bundle -o bundled-board.go board
//go:generate fyne bundle -o bundled-pieces.go pieces

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/notnil/chess"
)

type boardContainer struct {
	widget.BaseWidget

	objects []fyne.CanvasObject
	tapped  func()
}

func newBoardContainer(cells []fyne.CanvasObject, tap func()) *boardContainer {
	c := &boardContainer{objects: cells, tapped: tap}
	c.ExtendBaseWidget(c)
	return c
}

func (b *boardContainer) CreateRenderer() fyne.WidgetRenderer {
	return &boardRenderer{b: b, objects: b.objects}
}

func (b *boardContainer) Tapped(_ *fyne.PointEvent) {
	b.tapped()
}

type boardRenderer struct {
	b       *boardContainer
	objects []fyne.CanvasObject
}

func (b *boardRenderer) Destroy() {
}

func (b *boardRenderer) Layout(s fyne.Size) {
	cellEdge := cellSize(s)
	leftInset := (s.Width - cellEdge*8) / 2
	cellSize := fyne.NewSize(cellEdge, cellEdge)
	i := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			b.objects[i].Resize(cellSize)
			b.objects[i].Move(fyne.NewPos(
				leftInset+(float32(x)*cellEdge), float32(y)*cellEdge))

			i++
		}
	}
}

func (b *boardRenderer) MinSize()fyne.Size {
	edge := theme.IconInlineSize() * 8
	return fyne.NewSize(edge, edge)
}

func (b *boardRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *boardRenderer) Refresh() {
}

func cellSize(s fyne.Size) float32 {
	smallEdge := s.Width
	if s.Height < s.Width {
		smallEdge = s.Height
	}

	return smallEdge / 8
}

func resourceForPiece(p chess.Piece) fyne.Resource {
	switch p.Color() {
	case chess.Black:
		switch p.Type() {
		case chess.Pawn:
			return resourceBlackPawnSvg
		case chess.Rook:
			return resourceBlackRookSvg
		case chess.Knight:
			return resourceBlackKnightSvg
		case chess.Bishop:
			return resourceBlackBishopSvg
		case chess.Queen:
			return resourceBlackQueenSvg
		case chess.King:
			return resourceBlackKingSvg
		}
	case chess.White:
		switch p.Type() {
		case chess.Pawn:
			return resourceWhitePawnSvg
		case chess.Rook:
			return resourceWhiteRookSvg
		case chess.Knight:
			return resourceWhiteKnightSvg
		case chess.Bishop:
			return resourceWhiteBishopSvg
		case chess.Queen:
			return resourceWhiteQueenSvg
		case chess.King:
			return resourceWhiteKingSvg
		}
	}
	return nil
}
