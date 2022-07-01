package main

import (
	"encoding/json"
	"log"
	"os"

	goose "github.com/advancedlogic/GoOse"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cl := feeds.FeedCollectionsLoader{}

	mainFeeds, err := cl.LoadMain()
	if err != nil {
		return err
	}

	regionalFeeds, err := cl.LoadRegional()
	if err != nil {
		return err
	}

	fl := feeds.FeedRSSLoader{}

	m := map[string][]feeds.FeedItem{}

	for _, feed := range append(mainFeeds, regionalFeeds...) {
		rss, err := fl.Feed(feed.URL)
		if err != nil {
			return err
		}
		m[feed.Title] = (*rss).Channel.Items
	}

	bytes, _ := json.Marshal(m)
	_ = os.WriteFile("output.json", bytes, 0777)

	_ = goose.New()
	return nil
}
