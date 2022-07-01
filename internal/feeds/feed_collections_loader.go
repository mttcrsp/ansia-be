package feeds

import (
	"encoding/json"
	"os"
)

type FeedCollectionsLoader struct{}

func (l *FeedCollectionsLoader) LoadMain() ([]Feed, error) {
	return l.load("./assets/main-feeds.json")
}

func (l *FeedCollectionsLoader) LoadRegional() ([]Feed, error) {
	return l.load("./assets/regional-feeds.json")
}

func (l *FeedCollectionsLoader) load(path string) ([]Feed, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	feeds := []Feed{}
	err = json.Unmarshal(bytes, &feeds)
	return feeds, err
}
