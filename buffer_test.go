package main

import (
	"testing"
)

func TestSplitAlphabet(t *testing.T) {
	buf := newbuffer()
	alphabet := "abcdefg"
	for _, c := range alphabet {
		buf.insertChr(rune(c))
	}
	if l, r := buf.lines[buf.cursor.y].split(3); string(l) != "abc" || string(r) != "defg" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
	if l, r := buf.lines[buf.cursor.y].split(0); string(l) != "" || string(r) != "abcdefg" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
}

func TestSplitCJK(t *testing.T) {
	buf := newbuffer()
	cjk := "南無阿弥陀仏"
	for _, c := range cjk {
		buf.insertChr(rune(c))
	}
	if l, r := buf.lines[buf.cursor.y].split(3); string(l) != "南無阿" || string(r) != "弥陀仏" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
	if l, r := buf.lines[buf.cursor.y].split(0); string(l) != "" || string(r) != "南無阿弥陀仏" {
		t.Fatalf("failed test split function, l: %v, r: %v", string(l), string(r))
	}
}

func TestMoveCursor(t *testing.T) {
	buf := newbuffer()
	str1 := "123456"
	for _, c := range str1 {
		buf.insertChr(rune(c))
	}
	buf.moveCursor(Up)
	if buf.cursor.x != 6 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 6,0(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Down)
	if buf.cursor.x != 6 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 6,0(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Right)
	if buf.cursor.x != 6 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 6,0(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Left)
	if buf.cursor.x != 5 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 5,0(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	if buf.cursor.x != 0 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Right)
	if buf.cursor.x != 1 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 1,0(x, y), fact: %v", buf.cursor)
	}
	buf.lineFeed()
	//  1
	//  23456
	// ^ cursor
	if buf.cursor.x != 0 && buf.cursor.y != 1 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,1(x, y), fact: %v", buf.cursor)
	}
	buf.moveCursor(Up)
	if buf.cursor.x != 0 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0(x, y), fact: %v", buf.cursor)
	}

	buf.moveCursor(Down)
	buf.moveCursor(Right)
	buf.moveCursor(Right)
	buf.moveCursor(Up)
	if buf.cursor.x != 0 && buf.cursor.y != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0(x, y), fact: %v", buf.cursor)
	}
	str2 := "abcdefg"
	for _, c := range str2 {
		buf.insertChr(rune(c))
	}
	// abcdefg1
	// 23456
	buf.moveCursor(Down)
	if buf.cursor.x != 5 && buf.cursor.y != 1 {
		t.Fatalf("failed test moveCursor(Down), expected: 5,1, fact: %v", buf.cursor)
	}
}

func TestGetline(t *testing.T) {
	buf := newbuffer()
	if str, err := buf.getline(0); err != nil {
		if string(str) != "" {
			t.Fatalf("failed test getline(0), expected: '', fact: %v", string(str))
		}
	}

	str1 := "123456あいうえお"
	for _, c := range str1 {
		buf.insertChr(rune(c))
	}
	// only one line
	if str, err := buf.getline(0); err != nil {
		if string(str) != "123456あいうえお" {
			t.Fatalf("failed test getline(0), expected: 123456あいうえお, fact: %v", string(str))
		}
	}

	buf.lineFeed()
	str2 := "かきくけこさしすせそ"
	for _, c := range str2 {
		buf.insertChr(rune(c))
	}
	// have two line
	if str, err := buf.getline(0); err != nil {
		if string(str) != "123456あいうえお" {
			t.Fatalf("failed test getline(0), expected: 123456あいうえお, fact: %v", string(str))
		}
	}
	if str, err := buf.getline(1); err != nil {
		if string(str) != "かきくけこさしすせそ" {
			t.Fatalf("failed test getline(1), expected: かきくけこさしすせそ, fact: %v", string(str))
		}
	}
	if _, err := buf.getline(5); err == nil {
		t.Fatalf("failed test getline(1), expected: nil, fact: %v", err)
	}
}

func TestGetLastLine(t *testing.T) {
	buf := newbuffer()
	// バッファに何も書き込まれていない場合
	if buf.getlastlinenum() != 1 {
		t.Fatalf("failed test getlastlinenum(), expected: 1, fact: %v", buf.getlastlinenum())
	}
	str1 := "123456あいうえお"
	for _, c := range str1 {
		buf.insertChr(rune(c))
	}
	if buf.getlastlinenum() != 1 {
		t.Fatalf("failed test getlastlinenum(), expected: 1, fact: %v", buf.getlastlinenum())
	}

	buf.lineFeed()
	if buf.getlastlinenum() != 2 {
		t.Fatalf("failed test getlastlinenum(), expected: 2, fact: %v", buf.getlastlinenum())
	}

	str2 := "かきくけこさしすせそ"
	for _, c := range str2 {
		buf.insertChr(rune(c))
	}
	if buf.getlastlinenum() != 2 {
		t.Fatalf("failed test getlastlinenum(), expected: 2, fact: %v", buf.getlastlinenum())
	}
}

func TestDeleteLine(t *testing.T) {
	buf := newbuffer()
	if err := buf.deleteLine(1); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.getlastlinenum() != 1 {
		t.Fatalf("failed test deleteLine(), expected 0, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.y != 0 || newCur.x != 0 {
		t.Fatalf("failed test deleteLine(), expected 0,0 (x, y), fact: %v", newCur)
	}

	str1 := "12345あいうえ"
	for _, c := range str1 {
		buf.insertChr(rune(c))
		buf.lineFeed()
	}
	buf.insertChr(rune('お'))

	// カーソルよりも前にある行を削除
	if err := buf.deleteLine(1); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.getlastlinenum() != 9 {
		t.Fatalf("failed test deleteLine(), expected 9, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.y != 8 || newCur.x != 1 {
		t.Fatalf("failed test deleteLine(), expected 8,1(x, y), fact: %v", newCur)
	}

	buf.moveCursor(Up)
	buf.moveCursor(Up)
	buf.moveCursor(Up)

	// カーソルより後ろの行を削除
	if err := buf.deleteLine(7); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.getlastlinenum() != 8 {
		t.Fatalf("failed test deleteLine(), expected 8, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), face: %v", newCur)
	}

	// カーソル行を削除
	if err := buf.deleteLine(buf.cursor.y); err != nil {
		t.Fatalf("failed test deleteLine(), expected nil, fact: %v", err)
	}
	if buf.getlastlinenum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), face: %v", newCur)
	}

	// 存在しない行を削除 1
	if err := buf.deleteLine(100); err != ErrOutRange {
		t.Fatalf("failed test deleteLine(), expected ErrOutRange, fact: %v", err)
	}
	if buf.getlastlinenum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), fact: %v", newCur)
	}

	// 存在しない行を削除 2
	if err := buf.deleteLine(-1); err != ErrOutRange {
		t.Fatalf("failed test deleteLine(), expected ErrOutRange, fact: %v", err)
	}
	if buf.getlastlinenum() != 7 {
		t.Fatalf("failed test deleteLine(), expected 7, fact: %v", buf.getlastlinenum())
	}
	if newCur := buf.cursor; newCur.x != 1 || newCur.y != 5 {
		t.Fatalf("failed test deleteLine(), expected 1, 5(x, y), fact: %v", newCur)
	}
}

func TestUndoRedo(t *testing.T) {
	buf := newbuffer()
	buf.pushBufToUndoRedoBuffer()

	// undo buffer が空の時
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	l, err := buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	// redo buffer が空の時
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	alphabet := "abcdefg"
	for _, c := range alphabet {
		buf.insertChr(rune(c))
		buf.pushBufToUndoRedoBuffer()
	}
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdef" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdef', fact: %v", string(l))
	}

	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	buf.undo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "a" {
		t.Fatalf("failed test undo(), expected buffer: 'a', fact: %v", string(l))
	}

	buf.undo()
	buf.pushBufToUndoRedoBuffer()

	// 全ての undo バッファを引き出した状態
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "" {
		t.Fatalf("failed test undo(), expected buffer: '', fact: %v", string(l))
	}

	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "a" {
		t.Fatalf("failed test undo(), expected buffer: 'a', fact: %v", string(l))
	}

	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdefg" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdefg', fact: %v", string(l))
	}

	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	buf.redo()
	buf.pushBufToUndoRedoBuffer()
	l, err = buf.getline(0)
	if err != nil {
		t.Fatalf("failed test getline(0), expected: nil, fact: %v", err)
	}
	if string(l) != "abcdefg" {
		t.Fatalf("failed test undo(), expected buffer: 'abcdefg', fact: %v", string(l))
	}
}
