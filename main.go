package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var engine Chessplayer = RANDOM

func main() {
	if checkEngine() {
		engine = UCI
	}
	chessApp := app.NewWithID("xyz.andy.chess")
	win := chessApp.NewWindow("Chess")
	ui := newUI(win)
	game := newGame()
	game.initGame([2]Chessplayer{HUMAN, engine}, ui)
	game.LoadFromPreferences(chessApp)
	win.SetContent(game.ui.makeUI(game))
	win.Resize(fyne.NewSize(480, 480+theme.IconInlineSize()*2+theme.Padding()))
	win.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New Game", func() {
				dialogNewGame(&win, func(playerWhite, playerBlack Chessplayer) {
					game.Stop()
					game.initGame([2]Chessplayer{playerWhite, playerBlack}, ui)
					game.ui.refreshGrid(game.cgame)
					chessApp.Preferences().SetString(preferenceKeyCurrent, game.cgame.FEN())
					game.Play()
				})
			})),
	))
	game.Play()
	win.ShowAndRun()
}
