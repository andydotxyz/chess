package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var engine playerType = RANDOM

func main() {
	if checkEngine() {
		engine = UCI
	}
	chessApp := app.NewWithID("xyz.andy.chess")
	win := chessApp.NewWindow("Chess")
	ui := newUI(win)
	game := NewGame()
	game.InitGame([2]playerType{HUMAN, engine}, ui)
	win.SetContent(game.ui.makeUI(game))
	win.Resize(fyne.NewSize(480, 480+theme.IconInlineSize()*2+theme.Padding()))
	win.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New Game", func() {
				dialogNewGame(&win, func(playerWhite, playerBlack playerType) {
					game.Stop()
					game.InitGame([2]playerType{playerWhite, playerBlack}, ui)
					game.ui.refreshGrid(game.cgame)
					chessApp.Preferences().SetString(PREFERENCE_KEY_CURRENT, game.marshall())
					game.Play()
				})
			})),
	))

	cur := fyne.CurrentApp().Preferences().String(PREFERENCE_KEY_CURRENT)
	game.loadGame(cur, ui)
	game.ui.refreshGrid(game.cgame)
	game.Play()
	win.ShowAndRun()
}
