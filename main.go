package main

import (
	"fmt"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	buf := newbuffer()
	fmt.Print(len(os.Args))
	if len(os.Args) > 1 {
		buf.filename = os.Args[1]
	}
	if err := startUp(); err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()

	if buf.filename == "" {
		buf.lines = []*line{&line{[]rune{}}}
		buf.filename = "newfile.txt"
	} else {
		file, err := os.Open(buf.filename)
		if err != nil {
			log.Fatal(err)
		}
		buf.readFileToBuf(file)
	}

	//mode := Move

	buf.draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				buf.lineFeed()
			// mac delete-key is this
			case termbox.KeyCtrlH:
				fallthrough
			case termbox.KeyBackspace2:
				buf.backSpace()
			case termbox.KeyArrowUp:
				buf.moveCursor(Up)
			case termbox.KeyArrowDown:
				buf.moveCursor(Down)
			case termbox.KeyArrowLeft:
				buf.moveCursor(Left)
			case termbox.KeyArrowRight:
				buf.moveCursor(Right)
			case termbox.KeyCtrlZ:
				buf.undo()
			case termbox.KeyCtrlY:
				buf.redo()
			case termbox.KeyCtrlS:
				buf.writeBufToFile()
			case termbox.KeyCtrlC:
				break mainloop
			default:
				buf.insertChr(ev.Ch)
			}
		}

		// when entered any key, redraw buffer
		buf.draw()
	}
}

func startUp() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCursor(0, 0)
	return nil
}
