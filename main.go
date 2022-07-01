package main

import (
	"fmt"
	"log"

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
	rss, err := fl.Feed(mainFeeds[0].URL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", rss)

	_ = mainFeeds
	_ = regionalFeeds
	_ = goose.New()
	return nil
}
