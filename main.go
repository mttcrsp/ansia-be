package main

import (
	"fmt"
	"log"

	"github.com/mttcrsp/ansiabe/internal/articles"
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

	fl := feeds.RSSLoader{}
	rss, err := fl.Load(mainFeeds[0].URL)
	if err != nil {
		return err
	}

	article, err := articles.NewExtractor().Extract(rss.Channel.Items[0].Link)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", article)

	// regionalFeeds, err := cl.LoadRegional()
	// if err != nil {
	// 	return err
	// }
	// fl := feeds.RSSLoader{}

	// m := map[string][]feeds.Item{}

	// for _, feed := range append(mainFeeds, regionalFeeds...) {
	// 	rss, err := fl.Load(feed.URL)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	m[feed.Title] = (*rss).Channel.Items
	// }

	// bytes, _ := json.Marshal(m)
	// _ = os.WriteFile("output.json", bytes, 0777)

	return nil
}
