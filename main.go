package main

import (
	"log"
	"time"

	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	collectionsLoader := feeds.Loader{}

	mainFeeds, err := collectionsLoader.LoadMain()
	if err != nil {
		return err
	}

	regionalFeeds, err := collectionsLoader.LoadRegional()
	if err != nil {
		return err
	}

	loader := rss.Loader{}
	items := []rss.Item{}
	feeds := map[string][]rss.Item{}

	for _, feed := range append(mainFeeds, regionalFeeds...) {
		rss, err := loader.Load(feed.URL)
		if err != nil {
			return err
		}

		feedItems := (*rss).Channel.Items
		feeds[feed.Title] = feedItems
		items = append(items, feedItems...)
	}

	extractor := articles.NewExtractor()
	articles := map[string]articles.Article{}

	for _, item := range items {
		article, err := extractor.Extract(item.Link)
		if err != nil {
			return err
		}

		articles[item.Link] = article
		time.Sleep(time.Second / 2)
	}

	return nil
}
