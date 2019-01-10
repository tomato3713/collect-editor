package main

import "testing"

func TestEqual(t *testing.T) {
	mode := Move
	if b := mode.equal(Move); !b {
		t.Fatalf("failed test equal(), expected: true, fact: %v", b)
	}
	if b := mode.equal(Edit); b {
		t.Fatalf("failed test equal(), expected: false, fact: %v", b)
	}
	if b := mode.equal(Visual); b {
		t.Fatalf("failed test equal(), expected: false, fact: %v", b)
	}
	if b := mode.equal(Cmd); b {
		t.Fatalf("failed test equal(), expected: false, fact: %v", b)
	}
}

func TestString(t *testing.T) {
	mode := Move
	if l := mode.String(); l != "Move" {
		t.Fatalf("failed test String(), expected: 'Move', fact: %v", l)
	}
	mode = Edit
	if l := mode.String(); l != "Edit" {
		t.Fatalf("failed test String(), expected: 'Edit', fact: %v", l)
	}
	mode = Cmd
	if l := mode.String(); l != "Cmd" {
		t.Fatalf("failed test String(), expected: 'Cmd', fact: %v", l)
	}
	mode = Visual
	if l := mode.String(); l != "Visual" {
		t.Fatalf("failed test String(), expected: 'Visual', fact: %v", l)
	}
	var notInit Mode
	if l := notInit.String(); l != "Move" {
		t.Fatalf("failed test String(), expected: 'Move', fact: %v", l)
	}
}
