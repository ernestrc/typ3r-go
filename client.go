package typ3r

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Client represents a typ3r client
type Client struct {
	Config *Config
}

func (c *Client) request(method string, path string, data io.Reader) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response
	var buf []byte

	url := c.Config.serverURL + path
	client := &http.Client{}

	if req, err = http.NewRequest(method, url, data); err != nil {
		return
	}

	req.SetBasicAuth(c.Config.user, c.Config.token)

	if resp, err = client.Do(req); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200:
	case 403:
		err = fmt.Errorf(fmt.Sprintf("access denied to user %s with token %s", c.Config.user, c.Config.token))
		return
	default:
		err = fmt.Errorf(fmt.Sprintf("error: %+v\nrequest: %+v", resp, req))
		return
	}

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	return buf, nil
}

// ListNotes will return a list of the user notes
func (c *Client) ListNotes(offset int, limit int, query string) (notes Notes, err error) {
	var buf []byte
	path := fmt.Sprintf("/notes?limit=%d&offset=%d", limit, offset)

	if query != "" {
		path += "&query=" + query
	}

	if buf, err = c.request("GET", path, nil); err != nil {
		return
	}

	if err = json.Unmarshal(buf, &notes); err != nil {
		return
	}

	return notes, nil
}

func (c *Client) NewNote(text string) (*Note, error) {
	reader := strings.NewReader(fmt.Sprintf("{\"note\":\"%s\"}", text))

	var err error
	var buf []byte
	if buf, err = c.request("POST", "/notes", reader); err != nil {
		return nil, err
	}

	var note Note
	if err = json.Unmarshal(buf, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (c *Client) UpdateNote(id int, text string) error {
	var buf []byte
	var err error
	reader := strings.NewReader(text)
	path := fmt.Sprintf("/notes/%d", id)

	if buf, err = c.request("POST", path, reader); err != nil {
		return err
	}

	log.Println(string(buf))

	return nil
}
