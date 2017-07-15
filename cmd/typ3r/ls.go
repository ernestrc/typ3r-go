package main

import (
	"bytes"
	"fmt"
	"text/template"

	lessCfg "github.com/ernestrc/fractal/config"
	"github.com/ernestrc/fractal/less"
	typ3r "github.com/ernestrc/typ3r-go"
)

const fetchSize = 10

var (
	nextOffset int
	notes      []typ3r.Note
	search     string
)

func getNotes(client *typ3r.Client) (n int, err error) {
	var batch []typ3r.Note

	if nextOffset == -1 {
		return 0, nil
	}

	if batch, err = client.ListNotes(nextOffset, fetchSize, search); err != nil {
		return 0, err
	}

	n = len(batch)

	for _, n := range batch {
		notes = append(notes, n)
	}

	if nlen := len(batch); nlen == 0 {
		nextOffset = -1
	} else {
		nextOffset += nlen
	}

	return
}

func toString() string {
	var err error
	var buf bytes.Buffer
	var parsed *template.Template

	termwidth := less.Width()
	separator := ""

	for i := 0; i < termwidth; i++ {
		separator += "-"
	}

	parsed, err = template.New("notes").Parse(
		fmt.Sprintf("{{range $note := .}}%[1]s\n{{ $note.Text }}\n\n{{end}}", separator))

	if err != nil {
		panic(err)
	}

	if err = parsed.Execute(&buf, notes); err != nil {
		panic(err)
	}

	return buf.String()
}

func ls(client *typ3r.Client) error {
	var err error
	var n int

	if n, err = getNotes(client); err != nil {
		return err
	}

	cfg := lessCfg.New()
	cfg.Wrap = true

	if err = less.Init(cfg, ""); err != nil {
		return err
	}

	defer less.Close()

	for {
		if n != 0 {
			less.Content(toString())
		}

		less.Message(fmt.Sprintf("%d notes", len(notes)))

		ev := less.PollEvent()
		switch ev.Type {
		case less.EOF:
			n, err = getNotes(client)
		case less.Exit:
			return nil
		case less.Search:
			nextOffset = 0
			notes = notes[:0]
			search = string(ev.Data)
			less.Message(fmt.Sprintf("searching %s", search))
			n, err = getNotes(client)
		case less.Error:
			err = ev.Err
		}

		if err != nil {
			return err
		}
	}
}
