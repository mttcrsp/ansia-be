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
	l := feeds.FeedCollectionsLoader{}

	mainFeeds, err := l.LoadMain()
	if err != nil {
		return err
	}

	regionalFeeds, err := l.LoadRegional()
	if err != nil {
		return err
	}

	fmt.Println(mainFeeds)
	fmt.Println(regionalFeeds)
	_ = goose.New()
	return nil
}
