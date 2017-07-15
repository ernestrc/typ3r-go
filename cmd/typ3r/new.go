package main

import (
	"fmt"

	"github.com/ernestrc/sensible/editor"
	typ3r "github.com/ernestrc/typ3r-go"
)

func newNote(client *typ3r.Client) error {
	var data string
	var note *typ3r.Note
	var err error

	if data, err = editor.EditTmp(""); err != nil {
		return err
	}
	if note, err = client.NewNote(data); err != nil {
		return err
	}

	fmt.Printf("created new note with id %s", note.ID)

	return nil
}
