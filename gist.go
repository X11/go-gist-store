package gstate

import "encoding/json"

type Gist struct {
	Files map[string]*GistFile `json:"files"`
}

type GistFile struct {
	Content string `json:"content"`
}

func NewGist(bytes []byte) (Gist, error) {
	gist := Gist{}
	if err := json.Unmarshal(bytes, &gist); err != nil {
		return gist, err
	}
	return gist, nil
}

func (g *Gist) GetFileContent(name string) (string, bool) {
	if val, ok := g.Files[name]; ok {
		return val.Content, ok
	}
	return "", false
}

func (g *Gist) SetFileContent(name string, content string) bool {
	if _, ok := g.Files[name]; ok {
		g.Files[name].Content = content
		return true
	}
	return false
}

func (g *Gist) Marshal() ([]byte, error) {
	val, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	return val, nil
}
