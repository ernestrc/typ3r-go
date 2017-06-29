package typ3r

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"time"
)

// Timestamp is an alias to provide unmarshalling
type Timestamp time.Time

// Snippet holds a script or snippet of code
type Snippet struct {
}

// Task is a bummer
type Task struct {
}

// Note is the basic entity in typ3r
type Note struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Card     string    `json:"card"`
	Ts       Timestamp `json:"ts"`
	Created  Timestamp `json:"created"`
	Text     string    `json:"text"`
	ShareID  string    `json:"shareid"`
	Visits   string    `json:"visits"`
	Tasks    []Task    `json:"tasks"`
	Snippets []Snippet `json:"snippets"`
}

// Notes is simply a slice of Note
type Notes []Note

func (ln Notes) String() string {
	var err error
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 4, ' ', 0)

	if _, err = fmt.Fprintln(w, "id\tcard\tvisits\ttasks\tsnippets\ttext\tupdated_at"); err != nil {
		panic(err)
	}

	for _, n := range ln {
		if _, err = fmt.Fprintln(w, n.tabs()); err != nil {
			panic(err)
		}
	}

	w.Flush()
	return buf.String()
}

func (ct *Timestamp) UnmarshalJSON(b []byte) (err error) {
	var t time.Time
	normalized := strings.Join(strings.Split(strings.Trim(string(b), "\""), " "), "T") + "Z"
	t, err = time.Parse(time.RFC3339, normalized)
	*ct = Timestamp(t)
	return
}

func (n Note) tabs() string {
	limit := int(math.Min(float64(len(n.Text)), 10))
	summary := "\"" + strings.Replace(n.Text[0:limit], "\n", " ", -1) + "\""
	return fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t%s\t%s",
		n.ID, n.Card, n.Visits, len(n.Tasks), len(n.Snippets),
		summary, time.Time(n.Ts).String())
}