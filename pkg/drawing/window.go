package drawing

import (
	"github.com/homedm/collect-editor/pkg/buffer"

	termbox "github.com/nsf/termbox-go"
)

type Window struct {
	// The coordinates of the upper left corner
	Coord Coordinal
	Size  Size
	Buf   *buffer.Buffer
}

type Coordinal struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

func (w Window) drawHorizon(sx int, sy int, ex int, ey int, fg termbox.Attribute, bg termbox.Attribute) {
	var ch rune
	for x := sx; x < ex; x++ {
		for y := sy; y < ey; y++ {
			termbox.SetCell(x, y, ch, fg, bg)
		}
	}
}

func (w Window) UpdateCursor() {
	x, y := w.Buf.GetCursor()
	termbox.SetCursor(x+w.Coord.X, y+w.Coord.Y)
	// TODO: ある地点に動いていたら、ウィンドウ全体を動かす
}

func (w Window) UpdateBufBody() {
	// Clear inside this window
	for y := w.Coord.Y; y < w.Coord.Y+w.Size.Height; y++ {
		for x := w.Coord.X; x < w.Coord.X+w.Size.Width; x++ {
			termbox.SetCell(y, x, rune(0), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	// TODO: Draw text inside this window
	// 描画するテキストの範囲を決定する.
	// Draw text Stage
	l := w.Buf.GetLastLineNum()
	for y := 0; y < l; y++ {
		line, err := w.Buf.GetLine(y)
		if err != nil {
			return
		}
		for x, r := range line {
			termbox.SetCell(x+w.Coord.X, y+w.Coord.Y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}
