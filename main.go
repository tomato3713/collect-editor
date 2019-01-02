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

	mode := Move

	buf.draw()

	// poll for keyboard events in another goroutine
	events := make(chan termbox.Event, 1000)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

mainloop:
	for {
		select {
		case ev := <-events:
			switch mode {
			case Move:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyArrowUp:
						buf.moveCursor(Up)
					case termbox.KeyArrowDown:
						buf.moveCursor(Down)
					case termbox.KeyArrowLeft:
						buf.moveCursor(Left)
					case termbox.KeyArrowRight:
						buf.moveCursor(Right)
					case termbox.KeyCtrlS:
						buf.writeBufToFile()
					case termbox.KeyCtrlC:
						break mainloop // 実行終了
					default:
					}
					switch ev.Ch {
					case 'i':
						mode = Edit
					case 'u':
						buf.undo()
					case 'r':
						buf.redo()
					default:
					}
				}
			case Edit:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyEnter:
						buf.lineFeed()
						// mac delete-key is this
					case termbox.KeyCtrlH:
						fallthrough
					case termbox.KeyBackspace2:
						buf.backSpace()
					case termbox.KeyCtrlZ:
						buf.undo()
					case termbox.KeyCtrlY:
						buf.redo()
					default:
						buf.insertChr(ev.Ch)
					}
				}
			case Visual:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					}
				}
			default:
			}
			// when entered any key, redraw buffer
			buf.draw()
		default:
			// Nothing
		}
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
