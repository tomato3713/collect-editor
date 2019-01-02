package main

// Mode defines editor mode
type Mode int

// Move is mode for cursor moving
// Insert is mode for text inserting
// Visual is mode for text selecting
const (
	Move Mode = iota
	Insert
	Visual
)

// Mode has Stringer Interface
func (m Mode) String() string {
	switch m {
	case Move:
		return "Move"
	case Insert:
		return "Insert"
	case Visual:
		return "Visual"
	default:
		return "Unknown"
	}
}
