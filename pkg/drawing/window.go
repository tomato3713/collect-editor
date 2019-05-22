package drawing

import (
	"github.com/homedm/collect-editor/pkg/buffer"

	termbox "github.com/nsf/termbox-go"
)

// Window is window struct
type Window struct {
	// The coordinates of the upper left corner
	Coord Coordinal
	Size  Size
	Buf   *buffer.Buffer
}

// Coordinal is coordinal struct
type Coordinal struct {
	X int
	Y int
}

// Size is size struct
type Size struct {
	Width  int
	Height int
}

func (w Window) UpdateCursor() {
	x, y := w.Buf.GetPos()
	termbox.SetCursor(x+w.Coord.X, y+w.Coord.Y)
}

// Focus is focus the window
func (w Window) Focus() {
	w.UpdateCursor()
}

func (w Window) checkInside(x int, y int) (ok bool) {
	ok = true
	topX := w.Coord.X
	topY := w.Coord.Y
	bottomX := topX + w.Size.Width
	bottomY := topY + w.Size.Height

	if topX > x || bottomX < x {
		ok = false
	}
	if topY > y || bottomY < y {
		ok = false
	}
	return
}

func (w Window) UpdateBufBody() {
	w.clear()

	// TODO: Draw text inside this window
	// Draw text Stage
	bufLen := w.Buf.GetLastLineNum()

	for y := 0; y < bufLen; y++ {
		line, err := w.Buf.GetLine(y)
		if err != nil {
			return
		}
		for x, r := range line {
			DrawChr(x+w.Coord.X, y+w.Coord.Y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func (w Window) clear() {
	// Clear inside this window
	for y := w.Coord.Y; y < w.Coord.Y+w.Size.Height; y++ {
		for x := w.Coord.X; x < w.Coord.X+w.Size.Width; x++ {
			termbox.SetCell(y, x, rune(0), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

// Draw is draw window
func (w Window) Draw() {
	w.UpdateBufBody()
	w.Buf.PushBufToUndoRedoBuffer()
}
