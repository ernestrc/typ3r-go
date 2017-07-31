package main

import (
	"container/list"
	"fmt"

	"github.com/ernestrc/fractal"
	"github.com/ernestrc/fractal/component"
	"github.com/ernestrc/fractal/writer"
	typ3r "github.com/ernestrc/typ3r-go"
	termbox "github.com/nsf/termbox-go"
)

const selectColor = fractal.ColorRed
const resfg, resbg = fractal.AttrReverse, fractal.AttrReverse
const fg, bg = fractal.ColorDefault, fractal.ColorDefault
const fetchSize = 10
const tabspaces = 8
const wrap = false

type consoleMode uint8

const (
	notesMode consoleMode = iota
	searchMode
	editMode
)

var (
	width, height int
	w             writer.TermboxWriter
	wm            *component.TileManager
	cm            *component.TileManager
	notes         *component.List
	notesFrame    *component.Frame
	notesTile     *component.Tile
	noteScroll    *component.Scroll
	noteFrame     *component.Frame
	noteTile      *component.Tile
	noteBuf       fractal.Buffer
	searchScroll  *component.Scroll
	searchTile    *component.Tile
	searchBuf     fractal.Buffer
	msgScroll     *component.Scroll
	msgTile       *component.Tile
	msgBuf        fractal.Buffer
	cursor        fractal.Coordinates
	cursorHidden  bool
	mode          consoleMode
	prevMode      consoleMode
	tp            *typ3r.Client
	nextOffset    int
	search        string
	selected      *list.Element
	selectIdx     int
)

func setPreviousMode() error {
	if prevMode == notesMode {
		return setListMode()
	}

	if prevMode == searchMode {
		return setSearchMode()
	}

	return setEditMode()
}

func setSearchMode() (err error) {
	prevMode = mode
	mode = searchMode
	cursorHidden = false
	cursor.X, cursor.Y = 1, height-1
	searchBuf.Reset()
	if _, err = searchBuf.WriteRune('/'); err != nil {
		return
	}

	return
}

func setListMode() (err error) {
	mode = notesMode
	noteFrame.Fg = fractal.ColorDefault
	notesFrame.Fg = selectColor
	cursorHidden = true
	paintSelection(selectColor, fractal.AttrReverse)
	return nil

}

func setEditMode() (err error) {
	mode = editMode
	notesFrame.Fg = fractal.ColorDefault
	noteFrame.Fg = selectColor
	cursorHidden = false
	paintSelection(fractal.AttrReverse, fractal.AttrReverse)
	// TODO
	cursor.X, cursor.Y = 1, 1
	return nil
}

func doGetNotes() (n int, err error) {
	var batch []typ3r.Note

	if nextOffset == -1 {
		return 0, nil
	}

	if _, err = msgBuf.WriteString(fmt.Sprintf("searching %s, offset: %d", search, nextOffset)); err != nil {
		return
	}

	if batch, err = tp.ListNotes(nextOffset, fetchSize, search); err != nil {
		return 0, err
	}

	if n = len(batch); n == 0 {
		nextOffset = -1
	} else {
		for _, note := range batch {
			n := newFractalNote(&note)
			notes.PushBack(n)
		}
		nextOffset += n
	}

	return
}

func getNotes() (err error) {
	var n int
	for selectIdx == maxSelectIdx() && !notes.CanSeekDown() {
		n, err = doGetNotes()
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
	}

	return nil
}

func searchNotes(search string) (err error) {
	nextOffset = 0
	notes.Reset()
	msgBuf.Reset()

	if err = getNotes(); err != nil {
		return
	}

	msgBuf.Reset()
	if _, err = msgBuf.WriteString(fmt.Sprintf("%d notes", notes.Len())); err != nil {
		return
	}

	return nil
}

func getSelected() *Note {
	if selected == nil {
		return nil
	}
	return selected.Value.(*Note)
}

func writeEditBuffer() (err error) {
	noteBuf.Reset()
	if comp := getSelected(); comp != nil {
		_, err = noteBuf.WriteString(comp.note.Text)
	}
	return
}

func paintSelection(fg, bg fractal.Attribute) {
	if selected != nil {
		if comp := getSelected(); comp != nil {
			comp.Fg, comp.Bg = fg, bg
		}
	}
}

func resetSelection() {
	if selected != nil {
		if comp := getSelected(); comp != nil {
			comp.Fg, comp.Bg = fg, bg
		}
	}
}

func selectPrevNote() error {
	cursor.Y--
	resetSelection()
	if selected == nil {
		selectIdx = 0
		selected = notes.Front()
	} else if selected.Prev() != nil {
		selectIdx--
		selected = selected.Prev()
	}
	paintSelection(selectColor, fractal.AttrReverse)
	return writeEditBuffer()
}

func selectNextNote() error {
	resetSelection()
	if selected == nil {
		selectIdx = 0
		selected = notes.Front()
	} else if selected.Next() != nil {
		selectIdx++
		selected = selected.Next()
	}
	paintSelection(selectColor, fractal.AttrReverse)
	return writeEditBuffer()
}

func editHandle(ev termbox.Event) (exit bool, err error) {
	switch ev.Type {
	case termbox.EventResize:
		// TODO set cursor to proper place
		// cursor.X = searchBuf.Len()
		// cursor.Y = height - 1
	case termbox.EventKey:
		switch ev.Key {

		case termbox.KeyBackspace:
			fallthrough
		case termbox.KeyBackspace2:
			if cursor.X > 1 {
				cursor.X--
				searchBuf.Truncate(searchBuf.Len() - 1)
			}

		case termbox.KeyEnter:
			// ignore?
		case termbox.KeyEsc:
			if err = setListMode(); err != nil {
				return
			}

		default:
			// cursor.X++
			// if _, err = searchBuf.WriteRune(ev.Ch); err != nil {
			// 	return
			// }
		}
	}

	return
}

func searchHandle(ev termbox.Event) (exit bool, err error) {
	switch ev.Type {
	case termbox.EventResize:
		cursor.X = searchBuf.Len()
		cursor.Y = height - 1
	case termbox.EventKey:
		switch ev.Key {

		case termbox.KeyBackspace:
			fallthrough
		case termbox.KeyBackspace2:
			if cursor.X > 1 {
				cursor.X--
				searchBuf.Truncate(searchBuf.Len() - 1)
			}

		case termbox.KeyEnter:
			str := searchBuf.String()
			bytes := []byte(str)[1:]
			search = string(bytes)
			noteScroll.Search(search)
			if err = setPreviousMode(); err != nil {
				return
			}
			noteScroll.SeekNextResult()
			if err = searchNotes(search); err != nil {
				return
			}

		case termbox.KeyEsc:
			searchBuf.Reset()
			if err = setPreviousMode(); err != nil {
				return
			}

		default:
			cursor.X++
			if _, err = searchBuf.WriteRune(ev.Ch); err != nil {
				return
			}
		}
	}
	return
}

func maxSelectIdx() int {
	return notes.Height() + notes.Offset() - 1
}

func notesHandle(ev termbox.Event) (exit bool, err error) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEnter:
			setEditMode()
		case termbox.KeyEsc:
			return true, nil
		default:
			switch ev.Ch {
			case 'q':
				return true, nil
			case 'g':
				notes.SeekStart()
			case 'G':
				notes.SeekEnd()
			case 'j':
				if selectIdx == maxSelectIdx() {
					notes.SeekDown()
				}
				err = selectNextNote()
			case 'k':
				if selectIdx == notes.Offset() {
					notes.SeekUp()
				}
				err = selectPrevNote()
			case '/':
				err = setSearchMode()
			}
		}
	}
	return
}

func draw() (err error) {
	if cursorHidden {
		termbox.HideCursor()
	} else {
		termbox.SetCursor(cursor.X, cursor.Y)
	}
	if err = w.Clear(bg, fg); err != nil {
		return
	}
	if err = wm.Draw(&w); err != nil {
		return
	}
	if err = cm.Draw(&w); err != nil {
		return
	}
	if err = w.Flush(); err != nil {
		return
	}

	return
}

func resize() (err error) {
	if cm.Resize(width, 1); err != nil {
		return
	}
	if err = cm.Move(0, height-1); err != nil {
		return
	}
	return wm.Resize(width, height-1)
}

func poll() (err error) {
	var exit bool

	for !exit {
		if err = draw(); err != nil {
			return
		}
		ev := termbox.PollEvent()

		switch ev.Type {
		case termbox.EventError:
			err = ev.Err
		case termbox.EventResize:
			width, height = ev.Width, ev.Height
			if err = resize(); err != nil {
				break
			}
			fallthrough
		case termbox.EventKey:
			switch mode {
			case notesMode:
				exit, err = notesHandle(ev)
			case editMode:
				exit, err = editHandle(ev)
			case searchMode:
				exit, err = searchHandle(ev)
			}
		case termbox.EventMouse:
		case termbox.EventInterrupt:
		case termbox.EventRaw:
		case termbox.EventNone:
		}

		if err != nil {
			return
		}

		if err = getNotes(); err != nil {
			return
		}
	}

	return nil
}

// initialize console and run the event loop
func console(client *typ3r.Client) (err error) {
	if err = termbox.Init(); err != nil {
		return
	}

	defer termbox.Close()

	tp = client
	width, height = termbox.Size()

	notes = component.NewList(1, width, height, 0, 0)
	noteScroll = component.NewScroll(&noteBuf, 0, 0)
	noteScroll.Tabspaces = tabspaces
	noteScroll.Wrap = wrap

	if notesFrame, err = component.NewFrame(notes, 0, 0, 0, 0); err != nil {
		return
	}

	if noteFrame, err = component.NewFrame(noteScroll, 0, 0, 0, 0); err != nil {
		return
	}

	if wm, notesTile, err = component.NewTileManager(width, height-1, 0, 0, noteFrame); err != nil {
		return
	}

	if noteTile, err = wm.SplitHorizontal(notesTile, notesFrame); err != nil {
		return
	}

	searchScroll = component.NewScroll(&searchBuf, 0, 0)
	msgScroll = component.NewScroll(&msgBuf, 0, 0)

	if cm, searchTile, err = component.NewTileManager(width, 1, 0, height-1, searchScroll); err != nil {
		return
	}

	if msgTile, err = cm.SplitVertical(searchTile, msgScroll); err != nil {
		return
	}

	setListMode()

	return poll()
}
