// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"github.com/notnil/chess"
	"math/rand"
	"time"
)

type AgentRandom Agent

func NewAgentRandom() *AgentRandom {
	rand.Seed(time.Now().Unix())
	return &AgentRandom{}
}

func (a *AgentRandom) MakeMove(chessGame *chess.Game) *chess.Move {
	moves := chessGame.ValidMoves()
	return moves[rand.Intn(len(moves))]
}

func (a *AgentRandom) GetChannel() chan *chess.Move {
	return nil
}

func (a *AgentRandom) Stop() {
	return
}
