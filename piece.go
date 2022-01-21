package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/notnil/chess"
)

var (
	moveStart chess.Square = chess.NoSquare
)

type piece struct {
	widget.Icon
	g      *game
	u      *ui
	square chess.Square
}

func newPiece(u *ui, sq chess.Square, g *game) *piece {
	p := g.cgame.Position().Board().Piece(sq)

	ret := &piece{
		g:      g,
		u:      u,
		square: sq,
	}
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
	cell := p.u.grid.objects[off].(*fyne.Container)
	img := cell.Objects[2].(*piece)

	pos := cell.Position().Add(ev.Position)
	p.u.over.Move(pos.Subtract(fyne.NewPos(img.Size().Width/2, img.Size().Height/2)))
	p.u.over.Resize(img.Size())
	time.Sleep(1000)

	if img.Resource != nil {
		p.u.over.Resource = img.Resource
		p.u.over.Show()

		img.Resource = nil
		img.Refresh()
	}

	p.u.over.Refresh()
}

func (p *piece) DragEnd() {
	if moveStart != chess.NoSquare && p.square != moveStart {
		return // ignore drags if we are tapping
	}

	if p.u.start.Visible() {
		p.u.start.Hide()
		p.u.start.Refresh()
	}

	pos := p.u.over.Position().Add(fyne.NewPos(p.u.over.Size().Width/2, p.u.over.Size().Height/2))
	sq := positionToSquare(pos, p.u.grid.Size())

	var color int
	if p.g.cgame.Position().Turn() == chess.Black {
		color = 1
	}
	m := isValidMove(moveStart, sq, p.g.cgame)
	if m != nil {
		p.g.agents[color].GetChannel() <- m
	} else {
		off := squareToOffset(moveStart)
		cell := p.u.grid.objects[off].(*fyne.Container)
		pos2 := cell.Position()

		a := canvas.NewPositionAnimation(p.u.over.Position(), pos2, time.Millisecond*500, func(pos fyne.Position) {
			p.u.over.Move(pos)
			p.u.over.Refresh()
		})
		a.Start()
		time.Sleep(time.Millisecond * 550)
	}

	moveStart = chess.NoSquare
}

func (p *piece) Tapped(ev *fyne.PointEvent) {
	if moveStart == p.square {
		moveStart = chess.NoSquare
		p.u.start.Hide()
		p.u.start.Refresh()
		return
	}

	if moveStart == chess.NoSquare {
		if m := isValidMove(p.square, chess.NoSquare, p.g.cgame); m != nil {
			moveStart = p.square
			p.u.start.FillColor = okBGColor
			p.u.start.StrokeColor = okColor
		} else {
			p.u.start.FillColor = notOKBGColor
			p.u.start.StrokeColor = notOKColor
		}

		off := squareToOffset(p.square)
		cell := p.u.grid.objects[off].(*fyne.Container)

		p.u.start.Move(cell.Position())
		p.u.start.Resize(cell.Size())
		p.u.start.Refresh()
		p.u.start.Show()
		return
	}

	p.u.start.Hide()
	p.u.start.Refresh()

	off := squareToOffset(moveStart)
	cell := p.u.grid.objects[off].(*fyne.Container)

	var color int
	if p.g.cgame.Position().Turn() == chess.Black {
		color = 1
	}

	if m := isValidMove(moveStart, p.square, p.g.cgame); m != nil {
		moveStart = chess.NoSquare
		p.u.over.Move(cell.Position())
		p.g.agents[color].GetChannel() <- m
		return
	}

	moveStart = chess.NoSquare

	p.u.start.FillColor = notOKBGColor
	p.u.start.StrokeColor = notOKColor
	p.u.start.Move(cell.Position())
	p.u.start.Resize(cell.Size())
	p.u.start.Refresh()
	p.u.start.Show()

	go func() {
		time.Sleep(time.Millisecond * 500)
		p.u.start.Hide()
		p.u.start.Refresh()
	}()
}
