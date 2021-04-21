package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

var (
	okColor      = color.NRGBA{0, 0xff, 0, 0xff}
	okBGColor    = color.NRGBA{0, 0xff, 0, 0x28}
	notOKColor   = color.NRGBA{0xff, 0, 0, 0xff}
	notOKBGColor = color.NRGBA{0xff, 0, 0, 0x28}
)

type ui struct {
	grid  *fyne.Container
	over  *canvas.Image
	start *canvas.Rectangle
	win   fyne.Window

	game *chess.Game
	eng  *uci.Engine
}

func (u *ui) createGrid() *fyne.Container {
	var cells []fyne.CanvasObject

	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.NRGBA{0xF4, 0xE2, 0xB6, 0xFF})
			effect := canvas.NewImageFromResource(resourceOverlay1Png)
			effect.ScaleMode = canvas.ImageScaleFastest
			if x%2 == y%2 {
				bg.FillColor = color.RGBA{0x73, 0x50, 0x32, 0xFF}
				effect.Resource = resourceOverlay2Png
			}

			p := newPiece(u, chess.Square(x+y*8))
			cells = append(cells, container.NewMax(bg, effect, p))
		}
	}

	return container.New(&boardLayout{}, cells...)
}

func (u *ui) makeUI() fyne.CanvasObject {
	u.grid = u.createGrid()

	u.over = canvas.NewImageFromResource(nil)
	u.over.FillMode = canvas.ImageFillContain
	u.over.Hide()

	u.start = canvas.NewRectangle(color.Transparent)
	u.start.StrokeWidth = 4

	return container.NewMax(u.grid, container.NewWithoutLayout(u.start, u.over))
}


func (u *ui) refreshGrid() {
	y, x := 7, 0
	for _, cell := range u.grid.Objects {
		p := u.game.Position().Board().Piece(chess.Square(x + y*8))

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
