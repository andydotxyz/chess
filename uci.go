package main

import (
	"math/rand"
	"os/exec"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"

	"fyne.io/fyne/v2"
)

func loadOpponent() *uci.Engine {
	if _, err := exec.LookPath("stockfish"); err != nil {
		return nil
	}

	e, err := uci.New("stockfish") // you must have stockfish installed and on $PATH
	if err != nil {
		panic(err)
	}

	if err := e.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	return e
}

func playResponse(u *ui) {
	var m *chess.Move
	if u.eng != nil {
		cmdPos := uci.CmdPosition{Position: u.game.Position()}
		cmdGo := uci.CmdGo{MoveTime: time.Millisecond}
		if err := u.eng.Run(cmdPos, cmdGo); err != nil {
			panic(err)
		}

		m = u.eng.SearchResults().BestMove
	} else {
		m = randomResponse(u.game)
	}
	if m == nil {
		return // somehow end of game and we didn't notice?
	}

	off := squareToOffset(m.S1())
	cell := u.grid.objects[off].(*fyne.Container)

	u.over.Move(cell.Position())
	move(m, u.game, false, u, nil)
}

func randomResponse(game *chess.Game) *chess.Move {
	valid := game.ValidMoves()
	if len(valid) == 0 {
		return nil
	}

	return valid[rand.Intn(len(valid))]
}
