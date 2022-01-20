// Copyright © 2021-2022 Eugen Lindorfer, Austria

package main

import (
	"github.com/notnil/chess"
)

const (
	HUMAN playerType = iota
	RANDOM
	UCI
)

type AgentPlayer interface {
	MakeMove(*chess.Game) *chess.Move
	GetChannel() chan *chess.Move
	Stop()
}

type Agent struct {
	name string
}
