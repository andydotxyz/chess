// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Chessplayer int8

type game struct {
	gameId       uuid.UUID
	cgame        *chess.Game
	chessplayers [2]Chessplayer //TODO re-think names of attributes
	players      [2]AgentPlayer
	ui           *ui
	playing      bool
}

type gameSerial struct {
	PlayerWhite Chessplayer
	PlayerBlack Chessplayer
	FEN         string
	Uuid        uuid.UUID
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
	g.chessplayers = agents
	g.players = players
	g.ui = ui
	g.playing = true

	_ = g.ui.blackTurn.Set(false)
	_ = g.ui.outcome.Set(string(chess.NoOutcome))
}

func (g *game) loadGame(s string, ui *ui) {
	bytes := []byte(s)
	var gSerial gameSerial
	json.Unmarshal(bytes, &gSerial)
	agents := [2]Chessplayer{gSerial.PlayerWhite, gSerial.PlayerBlack}
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
	g.gameId = gSerial.Uuid
	load, err := chess.FEN(gSerial.FEN)
	if err != nil {
		return
	}
	load(g.cgame)

	g.chessplayers = agents
	g.players = players
	g.ui = ui
	g.playing = true

	_ = g.ui.blackTurn.Set(g.cgame.Position().Turn() == chess.White)
	_ = g.ui.outcome.Set(string(g.cgame.Outcome()))

	g.ui.refreshGrid(g.cgame)
}

func (g *game) marshall() string {
	gameSerial := &gameSerial{
		PlayerWhite: g.chessplayers[0],
		PlayerBlack: g.chessplayers[0],
		Uuid:        g.gameId,
		FEN:         g.cgame.FEN(),
	}
	b, _ := json.Marshal(gameSerial)
	return string(b)
}

func (g *game) Play() { // TODO think which methods need to be exported
	go func() {
		for {
			for color := 0; color < 2; color++ {
				if g.playing == false {
					return
				}

				if color == 0 && g.cgame.Position().Turn() == chess.Black {
					continue
				}

				m := g.players[color].MakeMove(g.cgame)
				move1(m, g.cgame, g.ui, g.chessplayers[color] != HUMAN)
				g.cgame.Move(m) // TODO handle the error
				move2(m, g.cgame, g.ui)
				if g.cgame.Outcome() != chess.NoOutcome {
					return
				}
				fyne.CurrentApp().Preferences().SetString(PREFERENCE_KEY_CURRENT, g.marshall())
				//TODO this should belong to the ui
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
