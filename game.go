// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type playerType int8

type game struct {
	gameId      uuid.UUID
	cgame       *chess.Game
	playerTypes [2]playerType
	agents      [2]AgentPlayer
	ui          *ui
	playing     bool
}

type gameSerial struct {
	PlayerWhite playerType
	PlayerBlack playerType
	FEN         string
	Id          uuid.UUID
}

func NewGame() *game {
	return &game{
		playing: false,
	}
}

func (g *game) InitGame(agents [2]playerType, ui *ui) {
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
	g.playerTypes = agents
	g.agents = players
	g.ui = ui
	g.playing = true

	_ = g.ui.blackTurn.Set(false)
	_ = g.ui.outcome.Set(string(chess.NoOutcome))

	//g.ui.refreshGrid(g.cgame)
}

func (g *game) loadGame(s string, ui *ui) {
	bytes := []byte(s)
	var gSerial gameSerial
	json.Unmarshal(bytes, &gSerial)
	agents := [2]playerType{gSerial.PlayerWhite, gSerial.PlayerBlack}
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
	g.gameId = gSerial.Id
	load, err := chess.FEN(gSerial.FEN)
	if err != nil {
		return
	}
	load(g.cgame)

	g.playerTypes = agents
	g.agents = players
	g.ui = ui
	g.playing = true

	_ = g.ui.blackTurn.Set(g.cgame.Position().Turn() == chess.White)
	_ = g.ui.outcome.Set(string(g.cgame.Outcome()))

	//g.ui.refreshGrid(g.cgame)
}

func (g *game) marshall() string {
	gameSerial := &gameSerial{
		PlayerWhite: g.playerTypes[0],
		PlayerBlack: g.playerTypes[0],
		Id:          g.gameId,
		FEN:         g.cgame.FEN(),
	}
	b, _ := json.Marshal(gameSerial)
	return string(b)
}

func (g *game) Play() {
	go func() {
		for {
			for color := 0; color < 2; color++ {
				if g.playing == false {
					return
				}

				if color == 0 && g.cgame.Position().Turn() == chess.Black {
					continue
				}

				m := g.agents[color].MakeMove(g.cgame)
				move1(m, g.cgame, g.ui, g.playerTypes[color] != HUMAN)
				err := g.cgame.Move(m)
				if err != nil {
					return
				}
				move2(m, g.cgame, g.ui)
				if g.cgame.Outcome() != chess.NoOutcome {
					return
				}
				fyne.CurrentApp().Preferences().SetString(PREFERENCE_KEY_CURRENT, g.marshall())
				//TODO this should belong to the ui
			}
		}
	}()
}

func (g *game) Stop() {
	g.playing = false
	if g.cgame.Outcome() == chess.NoOutcome {
		for color := 0; color < 2; color++ {
			g.agents[color].Stop()
		}
	}

}
