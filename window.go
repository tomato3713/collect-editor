package main

type window struct {
	coord   coordinal
	size    size
	visible bool
}

type coordinal struct {
	x int
	y int
}

type size struct {
	width  int
	height int
}
