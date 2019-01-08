package main

import (
	"container/list"
)

const (
	histLimit = 20
)

type cmdLine struct {
	histBuf list.List
	win     window
}

func newCmdLine() *cmdLine {
	c := new(cmdLine)
	c.histBuf.Init()
	c.histBuf.PushFront("")
	return c
}

func (c *cmdLine) addHist(cmd string) {
	if c.histBuf.Len() >= histLimit {
		// コマンド履歴がいっぱいになっている時は、最も古い履歴を削除してから、追加する。
		c.histBuf.Remove(c.histBuf.Back())
	}
	c.histBuf.Front().Value = cmd
	c.histBuf.PushFront("")
}
