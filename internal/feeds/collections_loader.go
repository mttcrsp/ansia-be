package feeds

import (
	"encoding/json"
	"os"
)

type CollectionsLoader struct{}

func (l *CollectionsLoader) LoadMain() ([]Feed, error) {
	return l.load("./assets/main-feeds.json")
}

func (l *CollectionsLoader) LoadRegional() ([]Feed, error) {
	return l.load("./assets/regional-feeds.json")
}

func (l *CollectionsLoader) load(path string) ([]Feed, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	feeds := []Feed{}
	err = json.Unmarshal(bytes, &feeds)
	return feeds, err
}
