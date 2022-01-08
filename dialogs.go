// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func dialogNewGame(mainWindow *fyne.Window, f func(playerWhite, playerBlack Chessplayer)) {
	var playerWhite Chessplayer = HUMAN
	var playerBlack Chessplayer = engine

	chooseWhite := &widget.RadioGroup{
		Options: []string{"Engine", "Human"},
		OnChanged: func(s string) {
			if s == "Engine" {
				playerWhite = engine
			} else {
				playerWhite = HUMAN
			}
		},
		Horizontal: true,
		Selected:   "Human",
	}

	chooseBlack := &widget.RadioGroup{
		Options: []string{"Engine", "Human"},
		OnChanged: func(s string) {
			if s == "Engine" {
				playerBlack = engine
			} else {
				playerBlack = HUMAN
			}
		},
		Horizontal: true,
		Selected:   "Engine",
	}

	items := []*widget.FormItem{
		widget.NewFormItem("Player White: ", chooseWhite),
		widget.NewFormItem("Player Black: ", chooseBlack),
	}

	dialog.ShowForm("New Game", "Start Game", "Cancel", items, func(ok bool) {
		if ok {
			f(playerWhite, playerBlack)
		}
	}, *mainWindow)

}
