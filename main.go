package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/notnil/chess"
)

func main() {
	a := app.New()
	w := a.NewWindow("Chess")

	game := chess.NewGame()
	grid := createGrid(game.Position().Board())

	over := canvas.NewImageFromResource(nil)
	over.Hide()
	bg := canvas.NewRectangle(color.Gray{Y: 0x7A})
	w.SetContent(container.NewMax(bg, grid, container.NewWithoutLayout(over)))
	w.Resize(fyne.NewSize(480, 480))

	go func() {
		rand.Seed(time.Now().Unix())
		for game.Outcome() == chess.NoOutcome {
			time.Sleep(time.Millisecond * 500)
			valid := game.ValidMoves()
			m := valid[rand.Intn(len(valid))]

			move(m, game, grid, over)
		}
	}()
	w.ShowAndRun()
}

func createGrid(b *chess.Board) *fyne.Container {
	grid := container.NewGridWithColumns(8)

	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.Gray{Y: 0xD0})
			if x%2 == y%2 {
				bg.FillColor = color.Gray{Y: 0x40}
			}

			p := b.Piece(chess.Square(x + y*8))
			img := canvas.NewImageFromResource(resourceForPiece(p))
			img.FillMode = canvas.ImageFillContain
			grid.Add(container.NewMax(bg, img))
		}
	}

	return grid
}

func refreshGrid(grid *fyne.Container, b *chess.Board) {
	y, x := 7, 0
	for _, cell := range grid.Objects {
		p := b.Piece(chess.Square(x + y*8))

		img := cell.(*fyne.Container).Objects[1].(*canvas.Image)
		img.Resource = resourceForPiece(p)
		img.Refresh()

		x++
		if x == 8 {
			x = 0
			y--
		}
	}
}

func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
	off := squareToOffset(m.S1())
	cell := grid.Objects[off].(*fyne.Container)
	img := cell.Objects[1].(*canvas.Image)
	pos1 := cell.Position()

	over.Resource = img.Resource
	over.Move(pos1)
	over.Resize(img.Size())

	over.Show()
	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	cell = grid.Objects[off].(*fyne.Container)
	pos2 := cell.Position()

	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
		over.Move(p)
		over.Refresh()
	})
	a.Start()
	time.Sleep(time.Millisecond * 550)

	game.Move(m)
	refreshGrid(grid, game.Position().Board())
	over.Hide()
}

func squareToOffset(sq chess.Square) int {
	x := sq % 8
	y := 7 - ((sq - x) / 8)

	return int(x + y*8)
}
