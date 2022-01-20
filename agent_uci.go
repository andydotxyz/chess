// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"math/rand"
	"os/exec"
	"time"
)

type AgentUCI struct {
	Agent
	playing     bool
	engine      *uci.Engine
	timePerMove time.Duration // millisecondes
}

func NewAgentUCI(timePerMove time.Duration) *AgentUCI {
	rand.Seed(time.Now().Unix())
	engine, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	if err := engine.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	ret := &AgentUCI{
		playing:     true,
		engine:      engine,
		timePerMove: timePerMove,
	}

	return ret
}

func (a *AgentUCI) MakeMove(chessGame *chess.Game) *chess.Move {
	cmdPos := uci.CmdPosition{Position: chessGame.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Millisecond * a.timePerMove}
	if err := a.engine.Run(cmdPos, cmdGo); err != nil {
		panic(err)
	}

	if a.playing == false {
		return nil
	}

	return a.engine.SearchResults().BestMove
}

func (a *AgentUCI) GetChannel() chan *chess.Move {
	return nil
}

func (a *AgentUCI) Stop() {
	a.playing = false
	return
}

func checkEngine() bool {
	if _, err := exec.LookPath("stockfish"); err != nil {
		return false
	}
	return true
}
