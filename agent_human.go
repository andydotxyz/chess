// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"github.com/notnil/chess"
	"log"
	"strconv"
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
	log.Print("human playing")

	for {
		select {
		case m := <-a.msg:
			log.Println("detected move by " + strconv.FormatBool(a.white))
			time.Sleep(500 * time.Millisecond)
			return m //TODO we might consider a check for valid moves here
		default:
			if a.playing == false {
				log.Print("agent was terminated")
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
	log.Print("stopping")
	return
}
