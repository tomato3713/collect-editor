package drawing

import termbox "github.com/nsf/termbox-go"

// DrawRect is draw horizontal line
func DrawRect(sx int, sy int, ex int, ey int, fg termbox.Attribute, bg termbox.Attribute) {
	var ch rune
	for x := sx; x < ex; x++ {
		for y := sy; y < ey; y++ {
			termbox.SetCell(x, y, ch, fg, bg)
		}
	}
}

// DrawChr is draw a character
func DrawChr(x int, y int, r rune, fg termbox.Attribute, bg termbox.Attribute) {
	termbox.SetCell(x, y, r, fg, bg)
}
