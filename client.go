package typ3r

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TPClient represents a typ3r client
type TPClient struct {
	Config *Config
}

func (c *TPClient) request(method string, path string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response
	var buf []byte

	url := c.Config.serverURL + path
	client := &http.Client{}

	if req, err = http.NewRequest(method, url, nil); err != nil {
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
		err = fmt.Errorf(fmt.Sprintf("error: %v", resp))
		return
	}

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	return buf, nil
}

// ListNotes will return a list of the user notes
func (c *TPClient) ListNotes() (notes Notes, err error) {
	var buf []byte

	if buf, err = c.request("GET", "/notes"); err != nil {
		return
	}

	if err = json.Unmarshal(buf, &notes); err != nil {
		return
	}

	return notes, nil
}
