package feeds

import (
	"encoding/json"
	"os"
)

type Loader struct{}

func (l *Loader) LoadMain() ([]Feed, error) {
	return l.load("./assets/main.json")
}

func (l *Loader) LoadRegional() ([]Feed, error) {
	return l.load("./assets/regional.json")
}

func (l *Loader) load(path string) ([]Feed, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	feeds := []Feed{}
	err = json.Unmarshal(bytes, &feeds)
	return feeds, err
}
