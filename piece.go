package main

import (
	"log"
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
	log.Println("dragged")

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

	if img.Resource != nil {
		p.u.over.Resource = img.Resource
		p.u.over.Show()

		img.Resource = nil
		img.Refresh()
	}

	p.u.over.Refresh()
}

func (p *piece) DragEnd() {
	log.Println("drag end")
	if moveStart != chess.NoSquare && p.square != moveStart {
		return // ignore drags if we are tapping
	}

	if p.u.start.Visible() {
		p.u.start.Hide()
		p.u.start.Refresh()
	}

	pos := p.u.over.Position().Add(fyne.NewPos(p.u.over.Size().Width/2, p.u.over.Size().Height/2))
	sq := positionToSquare(pos, p.u.grid.Size())

	log.Println("before isvalidmove")
	var color int
	if p.g.cgame.Position().Turn() == chess.Black {
		color = 1
	}
	m := isValidMove(moveStart, sq, p.g.cgame)
	log.Print(m)
	log.Print(p.g.chessplayers[color])
	if m != nil {
		//move(m, p.g.cgame, true, p.u)
		log.Println("before sending move")
		p.g.players[color].GetChannel() <- m
		log.Println("after sending move")

		/* go func() {
			time.Sleep(time.Second)
			//playResponse(p.u) TODO change this
		}()

		*/
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

		p.u.refreshGrid(p.g.cgame)
		p.u.over.Hide()
	}

	moveStart = chess.NoSquare
}

func (p *piece) Tapped(ev *fyne.PointEvent) {
	log.Println("tapped")
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
		//move(m, g.cgame, true, p.u) TODO change this
		p.g.players[color].GetChannel() <- m

		/*
			go func() {
				time.Sleep(time.Second / 2)
				//playResponse(p.u) TODO change this
			}()
		*/

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
