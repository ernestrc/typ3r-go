package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ararog/timeago"
	"github.com/ernestrc/less"
	typ3r "github.com/ernestrc/typ3r-go"
)

const fetchSize = 10

var (
	nextOffset int
	notes      []typ3r.Note
	search     string
)

func toRow(n *typ3r.Note) string {
	limit := int(math.Min(float64(len(n.Text)), 40))
	summary := "\"" + strings.Replace(n.Text[0:limit], "\n", " ", -1) + "\""

	got, err := timeago.TimeAgoFromNowWithTime(time.Time(n.Ts))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\t%s\t%d\t%d\t%s", summary, n.Visits,
		len(n.Tasks), len(n.Snippets), got)
}

func toTable(notes []typ3r.Note) string {
	var err error
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 4, ' ', 0)

	if _, err = fmt.Fprintln(w, "text\tvisits\ttasks\tsnippets\tupdated_at"); err != nil {
		panic(err)
	}

	for _, n := range notes {
		if _, err = fmt.Fprintln(w, toRow(&n)); err != nil {
			panic(err)
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}

	return buf.String()
}

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

func ls(client *typ3r.Client) error {
	var err error
	var n int

	if n, err = getNotes(client); err != nil {
		return err
	}

	if err = less.Init(nil, ""); err != nil {
		return err
	}

	defer less.Close()

	for {
		if n != 0 {
			less.Content(toTable(notes))
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
