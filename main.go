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
	cl := feeds.CollectionsLoader{}

	mainFeeds, err := cl.LoadMain()
	if err != nil {
		return err
	}

	regionalFeeds, err := cl.LoadRegional()
	if err != nil {
		return err
	}

	fl := feeds.RSSLoader{}

	m := map[string][]feeds.Item{}

	for _, feed := range append(mainFeeds, regionalFeeds...) {
		rss, err := fl.Load(feed.URL)
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
