package gstate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	GITHUB_API_URL = "https://api.github.com"
)

type GState struct {
	authentication string
	gistID         string
	gistURL        string
	gist           Gist
	filename       string
	fetched        bool
	rawState       string
}

// New GState instance. Requires the Gist ID, Filename and authentication string
func New(gistID string, filename string, authentication string) *GState {
	return &GState{
		authentication: "Basic " + base64.StdEncoding.EncodeToString([]byte(authentication)),
		gistID:         gistID,
		gistURL:        fmt.Sprintf("%s/gists/%s", GITHUB_API_URL, gistID),
		filename:       filename,
		fetched:        false,
	}
}

func (gs *GState) SetFetched(fetched bool) {
	gs.fetched = fetched
	if !fetched {
		gs.gist = Gist{}
		gs.rawState = ""
	}
}

// Get the state and marshal it to a referenced interface
func (gs *GState) Get(out interface{}) error {
	if gs.fetched == false {
		err := gs.fetch()
		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal([]byte(gs.rawState), out); err != nil {
		return err
	}
	return nil
}

func (gs *GState) Update(out interface{}) error {
	val, err := json.Marshal(out)
	if err != nil {
		return err
	}

	if ok := gs.gist.SetFileContent(gs.filename, string(val)); !ok {
		fmt.Errorf("Failed to update file '%s'", gs.filename)
	}

	val, err = gs.gist.Marshal()
	if err != nil {
		return err
	}

	resp, err := gs.githubRequest("PATCH", gs.gistURL, bytes.NewBuffer(val))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(string(body))
	}

	return nil
}

func (gs *GState) fetch() error {
	resp, err := gs.githubRequest("GET", gs.gistURL, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(string(body))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	gs.gist, err = NewGist(b)
	if err != nil {
		return err
	}

	var ok bool
	if gs.rawState, ok = gs.gist.GetFileContent(gs.filename); !ok {
		return fmt.Errorf("Could not find filename '%s' in gist", gs.filename)
	}
	gs.fetched = true

	return nil
}

func (gs *GState) githubRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", gs.authentication)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
