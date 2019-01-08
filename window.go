package main

import (
	termbox "github.com/nsf/termbox-go"
)

type window struct {
	// The coordinates of the upper left corner
	coord coordinal
	size  size
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
