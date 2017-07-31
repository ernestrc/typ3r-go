package main

import (
	"github.com/ernestrc/fractal"
	typ3r "github.com/ernestrc/typ3r-go"
)

type Note struct {
	note   *typ3r.Note
	width  int
	height int
	pos    fractal.Coordinates
	Fg     fractal.Attribute
	Bg     fractal.Attribute
}

func newFractalNote(note *typ3r.Note) (n *Note) {
	n = new(Note)
	n.note = note
	return
}

func (n *Note) Resize(width, height int) error {
	n.width, n.height = width, height
	return nil
}

func (n *Note) Move(x, y int) error {
	n.pos.X, n.pos.Y = x, y
	return nil
}

func (n *Note) Draw(w fractal.Writer) (err error) {
	runes := []rune(n.note.ID)
	for j, i, len := 0, n.pos.X, n.pos.X+len(runes); i < len; i, j = i+1, j+1 {
		if err = w.Write(i, n.pos.Y, runes[j], n.Fg, n.Bg); err != nil {
			return
		}
	}
	return nil
}
func (n *Note) Height() int {
	return n.height
}
func (n *Note) Width() int {
	return n.width
}
func (n *Note) Position() (int, int) {
	return n.pos.X, n.pos.Y
}
