package main

import (
	"fmt"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

var (
	mode Mode

	bufs       *Buffer
	cmdLineBuf *Buffer

	cmdLineWin CmdLineWin
	bufWins    BufferWin
)

func main() {
	if err := startUp(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	mode = Move
	bufs = newbuffer()
	cmdLineBuf = newbuffer()
	bufWins.buf = bufs
	cmdLineWin.buf = cmdLineBuf
	fmt.Print(len(os.Args))
	if len(os.Args) > 1 {
		bufs.filename = os.Args[1]
	}

	if bufs.filename == "" {
		bufs.filename = "newfile.txt"
	} else {
		file, err := os.Open(bufs.filename)
		if err != nil {
			log.Fatal(err)
		}
		bufs.readFileToBuf(file)
	}

	screenPaint()

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
						bufs.moveCursor(Up)
					case termbox.KeyArrowDown:
						bufs.moveCursor(Down)
					case termbox.KeyArrowLeft:
						bufs.moveCursor(Left)
					case termbox.KeyArrowRight:
						bufs.moveCursor(Right)
					case termbox.KeyCtrlS:
						bufs.writeBufToFile()
					case termbox.KeyCtrlC:
						break mainloop // 実行終了
					default:
					}
					switch ev.Ch {
					case ':':
						mode = Cmd
						cmdLineWin.focus()
					case 'i':
						mode = Edit
					case 'u':
						bufs.undo()
					case 'r':
						bufs.redo()
					default:
					}
				}
			case Edit:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyEnter:
						bufs.lineFeed()
						// mac delete-key is this
					case termbox.KeyCtrlH:
						fallthrough
					case termbox.KeyBackspace2:
						bufs.backSpace()
					case termbox.KeyCtrlZ:
						bufs.undo()
					case termbox.KeyCtrlR:
						bufs.redo()
					default:
						bufs.insertChr(ev.Ch)
					}
				}
			case Visual:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					}
				}
			case Cmd:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					default:
						cmdLineBuf.insertChr(ev.Ch)
					}
				}
			default:
			}
			// when entered any key, redraw buffer
			screenPaint()
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

	width, height := termbox.Size()
	// Set CmdLineWin default value
	cmdLineWin.coord.x = 0
	cmdLineWin.coord.y = height - cmdLineHeight
	cmdLineWin.size.width = width
	cmdLineWin.size.height = cmdLineHeight

	// Set bufWins default value
	bufWins.coord.x = 0
	bufWins.coord.y = 0
	bufWins.size.width = width
	bufWins.size.height = height - cmdLineHeight
	bufWins.stsLineHeight = 1
	return nil
}

func screenPaint() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	bufWins.draw()
	cmdLineWin.draw()

	if mode.equal(Cmd) {
		cmdLineWin.updateCursor()
	} else {
		bufWins.updateCursor()
	}
	termbox.Flush()
}
