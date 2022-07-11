package feeds

import (
	"encoding/json"
	"os"
)

type Loader struct{}

func (l *Loader) LoadCollections() (*Collections, error) {
	mainFeeds, err := l.load(
		collection{
			path: "./assets/main.json",
			slug: "principali",
		},
	)
	if err != nil {
		return nil, err
	}

	regionalFeeds, err := l.load(
		collection{
			path: "./assets/regional.json",
			slug: "regionali",
		},
	)
	if err != nil {
		return nil, err
	}

	mediaFeeds, err := l.load(
		collection{
			path: "./assets/media.json",
			slug: "media",
		},
	)
	if err != nil {
		return nil, err
	}

	return &Collections{
		Main:     mainFeeds,
		Regional: regionalFeeds,
		Media:    mediaFeeds,
	}, nil
}

func (l *Loader) load(collection collection) ([]Feed, error) {
	bytes, err := os.ReadFile(collection.path)
	if err != nil {
		return nil, err
	}

	feeds := []Feed{}
	if err = json.Unmarshal(bytes, &feeds); err != nil {
		return nil, err
	}

	for i := range feeds {
		feed := feeds[i]
		feed.CollectionSlug = collection.slug
		feeds[i] = feed
	}

	return feeds, nil
}

type collection struct {
	path string
	slug string
}
