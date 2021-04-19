package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

var (
	moveStart chess.Square = chess.NoSquare

	okColor      = color.NRGBA{0, 0xff, 0, 0xff}
	okBGColor    = color.NRGBA{0, 0xff, 0, 0x28}
	notOKColor   = color.NRGBA{0xff, 0, 0, 0xff}
	notOKBGColor = color.NRGBA{0xff, 0, 0, 0x28}
)

type piece struct {
	widget.Icon

	game   *chess.Game
	square chess.Square
}

func newPiece(g *chess.Game, sq chess.Square) *piece {
	p := g.Position().Board().Piece(sq)

	ret := &piece{game: g, square: sq}
	ret.ExtendBaseWidget(ret)

	ret.Resource = resourceForPiece(p)
	return ret
}

func (p *piece) Dragged(ev *fyne.DragEvent) {
	if moveStart != chess.NoSquare && p.square != moveStart {
		return // ignore drags if we are tapping
	}

	moveStart = p.square
	off := squareToOffset(p.square)
	cell := grid.Objects[off].(*fyne.Container)
	img := cell.Objects[2].(*piece)

	pos := cell.Position().Add(ev.Position)
	over.Move(pos.Subtract(fyne.NewPos(img.Size().Width/2, img.Size().Height/2)))
	over.Resize(img.Size())

	if img.Resource != nil {
		over.Resource = img.Resource
		over.Show()

		img.Resource = nil
		img.Refresh()
	}

	over.Refresh()
}

func (p *piece) DragEnd() {
	if moveStart != chess.NoSquare && p.square != moveStart {
		return // ignore drags if we are tapping
	}

	if start.Visible() {
		start.Hide()
		start.Refresh()
	}

	pos := over.Position().Add(fyne.NewPos(over.Size().Width/2, over.Size().Height/2))
	sq := positionToSquare(pos)

	if m := isValidMove(moveStart, sq, p.game); m != nil {
		move(m, p.game, grid, over)

		go func() {
			time.Sleep(time.Second)
			playResponse(p.game)
		}()
	} else {
		off := squareToOffset(moveStart)
		cell := grid.Objects[off].(*fyne.Container)
		pos2 := cell.Position()

		a := canvas.NewPositionAnimation(over.Position(), pos2, time.Millisecond*500, func(p fyne.Position) {
			over.Move(p)
			over.Refresh()
		})
		a.Start()
		time.Sleep(time.Millisecond * 550)

		refreshGrid(grid, p.game.Position().Board())
		over.Hide()
	}
	moveStart = chess.NoSquare
}

func (p *piece) Tapped(ev *fyne.PointEvent) {
	if moveStart == p.square {
		moveStart = chess.NoSquare
		start.Hide()
		start.Refresh()
		return
	}

	if moveStart == chess.NoSquare {
		if m := isValidMove(p.square, chess.NoSquare, p.game); m != nil {
			moveStart = p.square
			start.FillColor = okBGColor
			start.StrokeColor = okColor
		} else {
			start.FillColor = notOKBGColor
			start.StrokeColor = notOKColor
		}

		off := squareToOffset(p.square)
		cell := grid.Objects[off].(*fyne.Container)

		start.Move(cell.Position())
		start.Resize(cell.Size())
		start.Refresh()
		start.Show()
		return
	}

	start.Hide()
	start.Refresh()

	off := squareToOffset(moveStart)
	cell := grid.Objects[off].(*fyne.Container)

	if m := isValidMove(moveStart, p.square, p.game); m != nil {
		moveStart = chess.NoSquare
		over.Move(cell.Position())
		move(m, p.game, grid, over)

		go func() {
			time.Sleep(time.Second / 2)
			playResponse(p.game)
		}()
		return
	}

	moveStart = chess.NoSquare

	start.FillColor = notOKBGColor
	start.StrokeColor = notOKColor
	start.Move(cell.Position())
	start.Resize(cell.Size())
	start.Refresh()
	start.Show()

	go func() {
		time.Sleep(time.Millisecond * 500)
		start.Hide()
		start.Refresh()
	}()
}

func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()
	for _, m := range valid {
		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {
			return m
		}
	}

	return nil
}

func playResponse(game *chess.Game) {
	var m *chess.Move
	if eng != nil {
		cmdPos := uci.CmdPosition{Position: game.Position()}
		cmdGo := uci.CmdGo{MoveTime: time.Millisecond}
		if err := eng.Run(cmdPos, cmdGo); err != nil {
			panic(err)
		}

		m = eng.SearchResults().BestMove
	} else {
		m = randomResponse(game)
	}
	if m == nil {
		return // somehow end of game and we didn't notice?
	}

	off := squareToOffset(m.S1())
	cell := grid.Objects[off].(*fyne.Container)

	over.Move(cell.Position())
	move(m, game, grid, over)
}

func randomResponse(game *chess.Game) *chess.Move {
	valid := game.ValidMoves()
	if len(valid) == 0 {
		return nil
	}

	return valid[rand.Intn(len(valid))]
}
