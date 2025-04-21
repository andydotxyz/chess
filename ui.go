package main

import (
	"image/color"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	okColor      = color.NRGBA{0, 0xff, 0, 0xff}
	okBGColor    = color.NRGBA{0, 0xff, 0, 0x28}
	notOKColor   = color.NRGBA{0xff, 0, 0, 0xff}
	notOKBGColor = color.NRGBA{0xff, 0, 0, 0x28}
)

type ui struct {
	grid  *boardContainer
	over  *canvas.Image
	start *canvas.Rectangle
	win   fyne.Window

	blackTurn binding.Bool
	outcome   binding.String

	game *chess.Game
	eng  *uci.Engine
}

func newUI(win fyne.Window, game *chess.Game) *ui {
	u := &ui{win: win, game: game}
	u.blackTurn = binding.NewBool()
	u.outcome = binding.NewString()
	u.outcome.Set(string(chess.NoOutcome))

	return u
}

func (u *ui) createGrid() *boardContainer {
	var cells []fyne.CanvasObject

	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.NRGBA{0xF4, 0xE2, 0xB6, 0xFF})
			effect := canvas.NewImageFromResource(resourceOverlay1Png)
			effect.ScaleMode = canvas.ImageScaleFastest
			if x%2 == y%2 {
				bg.FillColor = color.RGBA{0x73, 0x50, 0x32, 0xFF}
				effect.Resource = resourceOverlay2Png
				effect.Image = nil
			}

			p := newPiece(u, chess.Square(x+y*8))
			cells = append(cells, container.NewMax(bg, effect, p))
		}
	}

	return newBoardContainer(cells, func() {
		moveStart = chess.NoSquare
		u.start.Hide()
		u.start.Refresh()
	})
}

func (u *ui) makeHeader() fyne.CanvasObject {
	whitePlays := widget.NewIcon(theme.NavigateBackIcon())
	blackPlays := widget.NewIcon(nil)
	status := widget.NewIcon(nil)

	u.blackTurn.AddListener(binding.NewDataListener(func() {
		if black, _ := u.blackTurn.Get(); black {
			whitePlays.SetResource(nil)
			blackPlays.SetResource(theme.NavigateNextIcon())
		} else {
			whitePlays.SetResource(theme.NavigateBackIcon())
			blackPlays.SetResource(nil)
		}
	}))
	u.outcome.AddListener(binding.NewDataListener(func() {
		outcome, _ := u.outcome.Get()
		switch outcome {
		case string(chess.NoOutcome):
			status.SetResource(nil)
		default:
			status.SetResource(theme.WarningIcon())
			blackPlays.SetResource(nil)
			whitePlays.SetResource(nil)
		}
	}))

	statusBG := canvas.NewRectangle(theme.BackgroundColor())
	statusBG.SetMinSize(fyne.NewSize(theme.IconInlineSize()*2, theme.IconInlineSize()*2))

	return container.NewHBox(layout.NewSpacer(),
		container.NewGridWithColumns(5,
			widget.NewIcon(resourceForPiece(chess.WhiteKing)),
			whitePlays,
			container.NewMax(statusBG, status),
			blackPlays,
			widget.NewIcon(resourceForPiece(chess.BlackKing)),
		),
		layout.NewSpacer())
}

func (u *ui) makeUI() fyne.CanvasObject {
	u.grid = u.createGrid()

	u.over = canvas.NewImageFromResource(nil)
	u.over.FillMode = canvas.ImageFillContain
	u.over.Hide()

	u.start = canvas.NewRectangle(color.Transparent)
	u.start.StrokeWidth = 4

	board := container.NewMax(u.grid, container.NewWithoutLayout(u.start, u.over))
	return container.NewBorder(u.makeHeader(), nil, nil, nil, board)
}

func (u *ui) refreshGrid() {
	y, x := 7, 0
	for _, cell := range u.grid.objects {
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
