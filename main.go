package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/server"
	"github.com/mttcrsp/ansiabe/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger := newLogger("core")
	extractor := articles.NewExtractor()
	feedsLoader := feeds.Loader{}
	rssLoader := rss.Loader{}
	store := store.Store{}

	collections, err := feedsLoader.LoadCollections()
	if err != nil {
		return err
	}

	feedsHandler := server.Feeds(
		server.FeedsVals{
			Collections: *collections,
		},
	)

	feedBySlugHandler := server.FeedBySlug(
		server.FeedBySlugVals{
			Collections: *collections,
		},
		server.FeedBySlugDeps{
			Store: store,
		},
	)

	c := make(chan string)

	go func() {
		gin.SetMode(gin.ReleaseMode)

		r := gin.Default()
		r.GET("/v1/feeds", feedsHandler)
		r.GET("/v1/feeds/:feed/items", feedBySlugHandler)
		r.Run()

		c <- "server did exit"
	}()

	go func() {
		for {
			for _, feed := range collections.All() {
				logger.Println("loading feed", feed.Slug())
				rssFeed, err := rssLoader.Load(feed.URL)
				if err != nil {
					logger.Printf("failed to load feed '%s': %s", feed.Slug(), err)
					time.Sleep(time.Second)
					continue
				}

				if err = store.InsertFeedItems(feed, *rssFeed); err != nil {
					logger.Printf("failed to insert items for feed '%s': %s", feed.Slug(), err)
					time.Sleep(time.Second)
					continue
				}

				for _, item := range rssFeed.Channel.Items {
					found, err := store.ArticleExists(item.ID())
					if err != nil {
						logger.Printf("failed to check article availability '%d': %s", item.ID(), err)
						time.Sleep(time.Second)
						continue
					}

					if found {
						time.Sleep(time.Second)
						continue
					}

					logger.Println("extracting article", item.Link)
					article, err := extractor.Extract(item.Link)
					if err != nil {
						logger.Printf("failed to extract article '%s': %s", item.Link, err)
						time.Sleep(time.Second)
						continue
					}

					if err = store.InsertArticle(item, *article); err != nil {
						logger.Printf("failed to insert article '%s: %s", item.Link, err)
						time.Sleep(time.Second)
						continue
					}
				}
			}

			time.Sleep(time.Minute)
		}
	}()

	return errors.New(<-c)
}

func newLogger(identifier string) *log.Logger {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)
	logger.SetPrefix(fmt.Sprintf("[%s] ", identifier))
	return logger
}
