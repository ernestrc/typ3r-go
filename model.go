package typ3r

import (
	"strings"
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

func (ct *Timestamp) UnmarshalJSON(b []byte) (err error) {
	var t time.Time
	normalized := strings.Join(strings.Split(strings.Trim(string(b), "\""), " "), "T") + "Z"
	t, err = time.Parse(time.RFC3339, normalized)
	*ct = Timestamp(t)
	return
}
