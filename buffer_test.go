package main

import (
	"testing"
)

func TestSplitAlphabet(t *testing.T) {
	buf := new(buffer)
	buf.lines = []*line{&line{[]rune{}}}
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
	buf := new(buffer)
	buf.lines = []*line{&line{[]rune{}}}
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
	buf := new(buffer)
	buf.lines = []*line{&line{[]rune{}}}
	str1 := "123456"
	for _, c := range str1 {
		buf.insertChr(rune(c))
	}
	buf.moveCursor(Up)
	if buf.cursor.y != 0 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 0,6, fact: %v,%v", buf.cursor.x, buf.cursor.y)
	}
	buf.moveCursor(Down)
	if buf.cursor.y != 0 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Up), expected: 0,6, fact: %v,%v", buf.cursor.x, buf.cursor.y)
	}
	buf.moveCursor(Right)
	if buf.cursor.y != 0 && buf.cursor.x != 6 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,6, fact: %v", buf.cursor.y)
	}
	buf.moveCursor(Left)
	if buf.cursor.y != 0 && buf.cursor.x != 5 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,5, fact: %v", buf.cursor.y)
	}
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	buf.moveCursor(Left)
	if buf.cursor.y != 0 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0, fact: %v", buf.cursor.y)
	}
	buf.moveCursor(Right)
	if buf.cursor.y != 0 && buf.cursor.x != 1 {
		t.Fatalf("failed test moveCursor(Down), expected: 1,0, fact: %v", buf.cursor.y)
	}
	buf.lineFeed()
	//  1
	//  23456
	// ^ cursor
	if buf.cursor.y != 1 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,1, fact: %v", buf.cursor.y)
	}
	buf.moveCursor(Up)
	if buf.cursor.y != 0 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0, fact: %v", buf.cursor.y)
	}

	buf.moveCursor(Down)
	buf.moveCursor(Right)
	buf.moveCursor(Right)
	buf.moveCursor(Up)
	if buf.cursor.y != 0 && buf.cursor.x != 0 {
		t.Fatalf("failed test moveCursor(Down), expected: 0,0, fact: %v", buf.cursor.y)
	}
	str2 := "abcdefg"
	for _, c := range str2 {
		buf.insertChr(rune(c))
	}
	// abcdefg1
	// 23456
	buf.moveCursor(Down)
	if buf.cursor.y != 1 && buf.cursor.x != 5 {
		t.Fatalf("failed test moveCursor(Down), expected: 5,1, fact: %v", buf.cursor.y)
	}
}

func TestGetline(t *testing.T) {
	buf := new(buffer)
	buf.lines = []*line{&line{[]rune{}}}
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
	buf := new(buffer)
	buf.lines = []*line{&line{[]rune{}}}
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
