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
	Top   int
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

func (w Window) getCursorPos() (x int, y int) {
	// ウィンドウ上の位置を計算する。
	// 左上 0, 0 とする
	x, y = w.Buf.GetPos()
	y = y - w.Top
	return
}

func (w Window) scroll() {
	_, y := w.getCursorPos()
	if y < 0 {
		w.Top--
	}
	if y >= w.Size.Height {
		w.Top++
	}
}

func (w Window) UpdateCursor() {
	x, y := w.getCursorPos()

	DrawCursor(x+w.Coord.X, y+w.Coord.Y)
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
	w.scroll()
	// TODO: Draw text inside this window
	// Draw text Stage
	for y := w.Coord.Y; y < w.Coord.Y+w.Size.Height; y++ {
		line, err := w.Buf.GetLine(w.Top + y - w.Coord.Y)
		if err != nil {
			return
		}
		for x, c := range line {
			DrawChr(x+w.Coord.X, y, c, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func (w Window) clear() {
	// Clear inside this window
	for y := w.Coord.Y; y < w.Coord.Y+w.Size.Height; y++ {
		for x := w.Coord.X; x < w.Coord.X+w.Size.Width; x++ {
			DrawChr(x, y, rune(0), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

// Draw is draw window
func (w Window) Draw() {
	w.UpdateBufBody()
	w.Buf.PushBufToUndoRedoBuffer()
}
