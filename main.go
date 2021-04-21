package main

import (
	"image/color"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/notnil/chess/uci"

	"github.com/notnil/chess"
)

const preferenceKeyCurrent = "current"

var (
	grid  *fyne.Container
	over  *canvas.Image
	start *canvas.Rectangle
	win   fyne.Window
	eng   *uci.Engine
)

func main() {
	a := app.NewWithID("xyz.andy.chess")
	win = a.NewWindow("Chess")

	game := chess.NewGame()
	loadGameFromPreference(game, a.Preferences())
	grid = createGrid(game)
	a.Preferences().AddChangeListener(func() {
		loadGameFromPreference(game, a.Preferences())
		refreshGrid(grid, game.Position().Board())
	})

	over = canvas.NewImageFromResource(nil)
	over.FillMode = canvas.ImageFillContain
	over.Hide()

	start = canvas.NewRectangle(color.Transparent)
	start.StrokeWidth = 4

	win.SetContent(container.NewMax(grid, container.NewWithoutLayout(start, over)))
	win.Resize(fyne.NewSize(480, 480))

	eng = loadOpponent()
	if eng != nil {
		defer eng.Close()
	} else {
		log.Println("Cound not find stockfish executable, using random player")
		rand.Seed(time.Now().Unix()) // random seed for random responses
	}
	win.ShowAndRun()
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

func loadOpponent() *uci.Engine {
	if _, err := exec.LookPath("stockfish"); err != nil {
		return nil
	}

	e, err := uci.New("stockfish") // you must have stockfish installed and on $PATH
	if err != nil {
		panic(err)
	}

	if err := e.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	return e
}

func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
	off := squareToOffset(m.S1())
	cell := grid.Objects[off].(*fyne.Container)
	img := cell.Objects[2].(*piece)

	over.Resource = resourceForPiece(game.Position().Board().Piece(m.S1()))
	over.Resize(img.Size())
	over.Refresh() // clear our old resource before showing

	over.Show()
	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	cell = grid.Objects[off].(*fyne.Container)
	pos2 := cell.Position()

	a := canvas.NewPositionAnimation(over.Position(), pos2, time.Millisecond*500, func(p fyne.Position) {
		over.Move(p)
		over.Refresh()
	})
	a.Start()
	time.Sleep(time.Millisecond * 550)

	game.Move(m)
	refreshGrid(grid, game.Position().Board())
	over.Hide()

	fyne.CurrentApp().Preferences().SetString(preferenceKeyCurrent, game.FEN())

	if game.Outcome() != chess.NoOutcome {
		result := "draw"
		switch game.Outcome().String() {
		case "1-0":
			result = "won"
		case "0-1":
			result = "lost"
		}

		fyne.CurrentApp().Preferences().SetString(preferenceKeyCurrent, "")
		dialog.ShowInformation("Game ended",
			"Game "+result+" because "+game.Method().String(), win)
	}
}

func positionToSquare(pos fyne.Position) chess.Square {
	var offX, offY = -1, -1
	for x := float32(0); x <= pos.X; x += grid.Size().Width / 8 {
		offX++
	}
	for y := float32(0); y <= pos.Y; y += grid.Size().Height / 8 {
		offY++
	}

	return chess.Square((7-offY)*8 + offX)
}

func squareToOffset(sq chess.Square) int {
	x := sq % 8
	y := 7 - ((sq - x) / 8)

	return int(x + y*8)
}
