package buffer

import (
	"testing"
)

func TestSplitAlphabet(t *testing.T) {
	buf := NewEmptyBuffer()
	alphabet := "abcdefg"
	for _, c := range alphabet {
		buf.InsertChr(rune(c))
	}
	if l, r := buf.lines[buf.pos.y].split(3); string(l) != "abc" || string(r) != "defg" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
	if l, r := buf.lines[buf.pos.y].split(0); string(l) != "" || string(r) != "abcdefg" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
}

func TestSplitCJK(t *testing.T) {
	buf := NewEmptyBuffer()
	cjk := "南無阿弥陀仏"
	for _, c := range cjk {
		buf.InsertChr(rune(c))
	}
	if l, r := buf.lines[buf.pos.y].split(3); string(l) != "南無阿" || string(r) != "弥陀仏" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
	if l, r := buf.lines[buf.pos.y].split(0); string(l) != "" || string(r) != "南無阿弥陀仏" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
}

func TestMoveCursor(t *testing.T) {
	buf := NewEmptyBuffer()
	str1 := "123456"
	for _, c := range str1 {
		buf.InsertChr(rune(c))
	}
	buf.MovePos(Up)
	if buf.pos.x != 6 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 6,0(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Down)
	if buf.pos.x != 6 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 6,0(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Right)
	if buf.pos.x != 6 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 6,0(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Left)
	if buf.pos.x != 5 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 5,0(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Left)
	buf.MovePos(Left)
	buf.MovePos(Left)
	buf.MovePos(Left)
	buf.MovePos(Left)
	buf.MovePos(Left)
	buf.MovePos(Left)
	if buf.pos.x != 0 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Right)
	if buf.pos.x != 1 && buf.pos.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 1,0(x, y), fact: %v", buf.pos)
	}
	buf.LineFeed()
	//  1
	//  23456
	// ^ cursor
	if buf.pos.x != 0 && buf.pos.y != 1 {
		t.Fatalf("failed test movepos(Down), expected: 0,1(x, y), fact: %v", buf.pos)
	}
	buf.MovePos(Up)
	if buf.pos.x != 0 && buf.pos.y != 0 {
		t.Fatalf("failed test movepos(Down), expected: 0,0(x, y), fact: %v", buf.pos)
	}

	buf.MovePos(Down)
	buf.MovePos(Right)
	buf.MovePos(Right)
	buf.MovePos(Up)
	if buf.pos.x != 0 && buf.pos.y != 0 {
		t.Fatalf("failed test movepos(Down), expected: 0,0(x, y), fact: %v", buf.pos)
	}
	str2 := "abcdefg"
	for _, c := range str2 {
		buf.InsertChr(rune(c))
	}
	// abcdefg1
	// 23456
	buf.MovePos(Down)
	if buf.pos.x != 5 && buf.pos.y != 1 {
		t.Fatalf("failed test movepos(Down), expected: 5,1, fact: %v", buf.pos)
	}
}

func TestGetline(t *testing.T) {
	buf := NewEmptyBuffer()
	if str, err := buf.GetLine(0); err != nil {
		if string(str) != "" {
			t.Fatalf("failed test getline(0), expected: '', fact: %v", string(str))
		}
	}

	str1 := "123456あいうえお"
	for _, c := range str1 {
		buf.InsertChr(rune(c))
	}
	// only one line
	if str, err := buf.GetLine(0); err != nil {
		if string(str) != "123456あいうえお" {
			t.Fatalf("failed test getline(0), expected: 123456あいうえお, fact: %v", string(str))
		}
	}

	buf.LineFeed()
	str2 := "かきくけこさしすせそ"
	for _, c := range str2 {
		buf.InsertChr(rune(c))
	}
	// have two line
	if str, err := buf.GetLine(0); err != nil {
		if string(str) != "123456あいうえお" {
			t.Fatalf("failed test getline(0), expected: 123456あいうえお, fact: %v", string(str))
		}
	}
	if str, err := buf.GetLine(1); err != nil {
		if string(str) != "かきくけこさしすせそ" {
			t.Fatalf("failed test getline(1), expected: かきくけこさしすせそ, fact: %v", string(str))
		}
	}
	if _, err := buf.GetLine(5); err == nil {
		t.Fatalf("failed test getline(1), expected: nil, fact: %v", err)
	}
}

func TestGetLastLine(t *testing.T) {
	buf := NewEmptyBuffer()
	// バッファに何も書き込まれていない場合
	if buf.GetLastLineNum() != 1 {
		t.Fatalf("failed test getlastlinenum(), expected: 1, fact: %v", buf.GetLastLineNum())
	}
	str1 := "123456あいうえお"
	for _, c := range str1 {
		buf.InsertChr(rune(c))
	}
	if buf.GetLastLineNum() != 1 {
		t.Fatalf("failed test getlastlinenum(), expected: 1, fact: %v", buf.GetLastLineNum())
	}

	buf.LineFeed()
	if buf.GetLastLineNum() != 2 {
		t.Fatalf("failed test getlastlinenum(), expected: 2, fact: %v", buf.GetLastLineNum())
	}

	str2 := "かきくけこさしすせそ"
	for _, c := range str2 {
		buf.InsertChr(rune(c))
	}
	if buf.GetLastLineNum() != 2 {
		t.Fatalf("failed test getlastlinenum(), expected: 2, fact: %v", buf.GetLastLineNum())
	}
}

func TestDeleteLine(t *testing.T) {
	buf := NewEmptyBuffer()
	if err := buf.DeleteLine(1); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.GetLastLineNum() != 1 {
		t.Fatalf("failed test deleteLine(), expected 0, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.y != 0 || newCur.x != 0 {
		t.Fatalf("failed test deleteLine(), expected 0,0 (x, y), fact: %v", newCur)
	}

	str1 := "12345あいうえ"
	for _, c := range str1 {
		buf.InsertChr(rune(c))
		buf.LineFeed()
	}
	buf.InsertChr(rune('お'))

	// カーソルよりも前にある行を削除
	if err := buf.DeleteLine(1); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.GetLastLineNum() != 9 {
		t.Fatalf("failed test deleteLine(), expected 9, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.y != 8 || newCur.x != 1 {
		t.Fatalf("failed test deleteLine(), expected 8,1(x, y), fact: %v", newCur)
	}

	buf.MovePos(Up)
	buf.MovePos(Up)
	buf.MovePos(Up)

	// カーソルより後ろの行を削除
	if err := buf.DeleteLine(7); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.GetLastLineNum() != 8 {
		t.Fatalf("failed test deleteLine(), expected 8, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), face: %v", newCur)
	}

	// カーソル行を削除
	if err := buf.DeleteLine(buf.pos.y); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.GetLastLineNum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), face: %v", newCur)
	}

	// 存在しない行を削除 1
	if err := buf.DeleteLine(100); err != ErrOutRange {
		t.Fatalf("failed test deleteLine(), expected ErrOutRange, fact: %v", err)
	}
	if buf.GetLastLineNum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), fact: %v", newCur)
	}

	// 存在しない行を削除 2
	if err := buf.DeleteLine(-1); err != ErrOutRange {
		t.Fatalf("failed test deleteLine(), expected ErrOutRange, fact: %v", err)
	}
	if buf.GetLastLineNum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.GetLastLineNum())
	}
	if newCur := buf.pos; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), fact: %v", newCur)
	}
}

func TestUndoRedo(t *testing.T) {
	buf := NewEmptyBuffer()
	buf.PushBufToUndoRedoBuffer()

	// undo buffer が空の時
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	l, err := buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	// redo buffer が空の時
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	alphabet := "abcdefg"
	for _, c := range alphabet {
		buf.InsertChr(rune(c))
		buf.PushBufToUndoRedoBuffer()
	}
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdef" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdef', fact: %v", string(l))
	}

	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	buf.Undo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "a" {
		t.Fatalf("failed test undo(), expected buffer: 'a', fact: %v", string(l))
	}

	buf.Undo()
	buf.PushBufToUndoRedoBuffer()

	// 全ての undo バッファを引き出した状態
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "a" {
		t.Fatalf("failed test undo(), expected buffer: 'a', fact: %v", string(l))
	}

	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdefg" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdefg', fact: %v", string(l))
	}

	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	buf.Redo()
	buf.PushBufToUndoRedoBuffer()
	l, err = buf.GetLine(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdefg" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdefg', fact: %v", string(l))
	}
}
