package main

import (
	"github.com/homedm/collect-editor/pkg/buffer"

	termbox "github.com/nsf/termbox-go"
)

type window struct {
	// The coordinates of the upper left corner
	coord coordinal
	size  size
	buf   *buffer.Buffer
}

type coordinal struct {
	x int
	y int
}

type size struct {
	width  int
	height int
}

func (w window) drawHorizon(sx int, sy int, ex int, ey int, fg termbox.Attribute, bg termbox.Attribute) {
	var ch rune
	for x := sx; x < ex; x++ {
		for y := sy; y < ey; y++ {
			termbox.SetCell(x, y, ch, fg, bg)
		}
	}
}

func (w window) updateCursor() {
	x, y := w.buf.GetCursor()
	termbox.SetCursor(x+w.coord.x, y+w.coord.y)
	// TODO: ある地点に動いていたら、ウィンドウ全体を動かす
}

func (w window) updateBufBody() {
	// Clear inside this window
	for y := w.coord.y; y < w.coord.y+w.size.height; y++ {
		for x := w.coord.x; x < w.coord.x+w.size.width; x++ {
			termbox.SetCell(y, x, rune(0), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	// TODO: Draw text inside this window
	// 描画するテキストの範囲を決定する.
	// Draw text Stage
	l := w.buf.GetLastLineNum()
	for y := 0; y < l; y++ {
		line, err := w.buf.GetLine(y)
		if err != nil {
			return
		}
		for x, r := range line {
			termbox.SetCell(x+w.coord.x, y+w.coord.y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}
