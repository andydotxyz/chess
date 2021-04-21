package main

import (
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"

	"github.com/notnil/chess"
)

func main() {
	a := app.NewWithID("xyz.andy.chess")
	win := a.NewWindow("Chess")

	game := chess.NewGame()
	u := &ui{win: win, game: game, eng: loadOpponent()}
	if u.eng != nil {
		defer u.eng.Close()
	} else {
		log.Println("Cound not find stockfish executable, using random player")
		rand.Seed(time.Now().Unix()) // random seed for random responses
	}

	loadGameFromPreference(game, a.Preferences())
	a.Preferences().AddChangeListener(func() {
		loadGameFromPreference(game, a.Preferences())
		u.refreshGrid()
	})

	win.SetContent(u.makeUI())
	win.Resize(fyne.NewSize(480, 480))

	win.ShowAndRun()
}

func move(m *chess.Move, game *chess.Game, u *ui) {
	off := squareToOffset(m.S1())
	cell := u.grid.Objects[off].(*fyne.Container)
	img := cell.Objects[2].(*piece)

	u.over.Resource = resourceForPiece(game.Position().Board().Piece(m.S1()))
	u.over.Resize(img.Size())
	u.over.Refresh() // clear our old resource before showing

	u.over.Show()
	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	cell = u.grid.Objects[off].(*fyne.Container)
	pos2 := cell.Position()

	a := canvas.NewPositionAnimation(u.over.Position(), pos2, time.Millisecond*500, func(p fyne.Position) {
		u.over.Move(p)
		u.over.Refresh()
	})
	a.Start()
	time.Sleep(time.Millisecond * 550)

	game.Move(m)
	u.refreshGrid()
	u.over.Hide()

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
			"Game "+result+" because "+game.Method().String(), u.win)
	}
}
