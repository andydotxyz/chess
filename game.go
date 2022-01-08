// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Chessplayer int8

type game struct {
	gameId       uuid.UUID
	cgame        *chess.Game
	chessplayers [2]Chessplayer
	players      [2]AgentPlayer
	ui           *ui
	playing      bool
}

func newGame() *game {
	return &game{
		playing: false,
	}
}

func (g *game) initGame(agents [2]Chessplayer, ui *ui) {
	var players [2]AgentPlayer

	for color := 0; color < 2; color++ {
		switch agents[color] {
		case HUMAN:
			players[color] = NewAgentHuman(color == 0)
		case RANDOM:
			players[color] = NewAgentRandom()
		case UCI:
			players[color] = NewAgentUCI(100)
		}
	}

	g.gameId = uuid.New()
	g.cgame = chess.NewGame()
	g.chessplayers = agents //TODO rethink naming of attributes
	g.players = players
	g.ui = ui
	g.playing = true

	_ = g.ui.blackTurn.Set(false)
	_ = g.ui.outcome.Set(string(chess.NoOutcome))
}

func (g *game) Play() {
	go func() {
		for {
			for color := 0; color < 2; color++ {
				if g.playing == false {
					return
				}

				m := g.players[color].MakeMove(g.cgame)
				move1(m, g.cgame, g.ui, g.chessplayers[color] != HUMAN)
				g.cgame.Move(m) // TODO handle the error
				move2(m, g.cgame, g.ui)
				if g.cgame.Outcome() != chess.NoOutcome {
					return
				}
			}
			//time.Sleep(200 * time.Millisecond)
		}
	}()
}

func (g *game) Stop() {
	if g.cgame.Outcome() == chess.NoOutcome {
		for color := 0; color < 2; color++ {
			g.players[color].Stop()
		}

		g.playing = false
	}

}

func (g *game) LoadFromPreferences(a fyne.App) {
	g.loadGameFromPreference(a.Preferences())
	a.Preferences().AddChangeListener(func() {
		g.loadGameFromPreference(a.Preferences())
		g.ui.refreshGrid(g.cgame)
	})
}

func (g *game) loadGameFromPreference(p fyne.Preferences) {
	cur := p.String(preferenceKeyCurrent)
	if cur == "" {
		return
	}

	load, err := chess.FEN(cur)
	if err != nil {
		return
	}
	load(g.cgame)
}
