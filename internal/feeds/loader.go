package feeds

import (
	"encoding/json"
	"os"
)

type Loader struct{}

func (l *Loader) LoadCollections() (*Collections, error) {
	mainFeeds, err := l.LoadMain()
	if err != nil {
		return nil, err
	}

	regionalFeeds, err := l.LoadRegional()
	if err != nil {
		return nil, err
	}

	mediaFeeds, err := l.LoadMedia()
	if err != nil {
		return nil, err
	}

	return &Collections{
		Main:     mainFeeds,
		Regional: regionalFeeds,
		Media:    mediaFeeds,
	}, nil
}

func (l *Loader) LoadMain() ([]Feed, error) {
	return l.load("./assets/main.json")
}

func (l *Loader) LoadRegional() ([]Feed, error) {
	return l.load("./assets/regional.json")
}

func (l *Loader) LoadMedia() ([]Feed, error) {
	return l.load("./assets/media.json")
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
