// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"github.com/notnil/chess"
	"time"
)

type AgentHuman struct {
	Agent
	msg       chan *chess.Move
	terminate chan bool
	white     bool
	playing   bool
}

func NewAgentHuman(white bool) *AgentHuman {
	return &AgentHuman{
		msg:       make(chan *chess.Move),
		terminate: make(chan bool),
		playing:   true,
	}
}

func (a *AgentHuman) MakeMove(chessGame *chess.Game) *chess.Move {
	for {
		select {
		case m := <-a.msg:
			time.Sleep(500 * time.Millisecond)
			return m //TODO we might consider a check for valid moves here
		default:
			if a.playing == false {
				return nil
			}

		}
	}
}

func (a *AgentHuman) GetChannel() chan *chess.Move {
	return a.msg
}

func (a *AgentHuman) Stop() {
	a.playing = false
	return
}
