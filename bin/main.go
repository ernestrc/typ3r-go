package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/ernestrc/less"
	"github.com/ernestrc/sensible/editor"
	typ3r "github.com/ernestrc/typ3r-go"
)

const fetchSize = 10

type typ3rCLI struct {
	client     *typ3r.Client
	nextOffset int
}

func (c typ3rCLI) WriteTo(w io.Writer) (n int64, err error) {
	var notes typ3r.Notes

	if c.nextOffset == -1 {
		return
	}

	if notes, err = c.client.ListNotes(c.nextOffset, fetchSize, ""); err != nil {
		return
	}

	var nlen int
	if nlen = len(notes); nlen == 0 {
		c.nextOffset = -1
		return
	}

	c.nextOffset += nlen

	var nn int
	if nn, err = w.Write([]byte(notes.String())); err != nil {
		return
	}

	n = int64(nn)

	return
}

func (c typ3rCLI) OnSearch(view *less.Handle, text string) error {
	return nil
}

func main() {
	var err error
	var args map[string]interface{}
	var config *typ3r.Config
	var view *less.Handle

	usage := `Typ3r Client

Usage:
  typ3r ls
  typ3r new
  typ3r -h | --help
  typ3r --version

Options:
  -h --help     Show this screen.`

	if args, err = docopt.Parse(usage, nil, true, "typ3r CLI", false); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if config, err = typ3r.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := typ3r.Client{Config: config}
	CLI := typ3rCLI{client: &client}

	if view = less.New(CLI, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if args["ls"].(bool) {
		err = view.Run()
	} else if args["new"].(bool) {
		var data string
		var note *typ3r.Note
		if data, err = editor.EditTmp(""); err != nil {
			goto exit
		}
		if note, err = client.NewNote(data); err != nil {
			goto exit
		}
		fmt.Printf("created new note with id %s", note.ID)
	} else {
		fmt.Println(usage)
	}

exit:
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
