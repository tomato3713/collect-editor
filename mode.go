package main

// Mode defines editor mode
type Mode int

// Move is mode for cursor moving
// Edit is mode for text editing
// Visual is mode for text selecting
const (
	Move Mode = iota
	Edit
	Visual
	Cmd
)

// Mode has Stringer Interface
func (m Mode) String() string {
	switch m {
	case Move:
		return "Move"
	case Edit:
		return "Edit"
	case Visual:
		return "Visual"
	case Cmd:
		return "Cmd"
	default:
		return "Unknown"
	}
}
