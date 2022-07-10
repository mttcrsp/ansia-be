package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/core"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/server"
	"github.com/mttcrsp/ansiabe/internal/store"
	"github.com/mttcrsp/ansiabe/internal/videojournal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	feedsLoader := feeds.Loader{}
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

	videojournalhandler := server.Videojournal(
		server.VideojournalDeps{
			Store: store,
		},
	)

	c := make(chan string)

	go func() {
		gin.SetMode(gin.ReleaseMode)

		r := gin.Default()
		r.GET("/v1/feeds", feedsHandler)
		r.GET("/v1/feeds/:feed/items", feedBySlugHandler)
		r.GET("/v1/videojournals", videojournalhandler)
		r.Run()

		c <- "server did exit"
	}()

	rssLoader := rss.Loader{}
	logger := newLogger("core")
	processors := []core.RSSProcessor{
		core.NewArticlesProcessor(
			*articles.NewExtractor(),
			store,
			*newLogger("articles"),
		),
		core.NewVideojournalProcessor(
			videojournal.Extractor{},
			store,
			*newLogger("videojournal"),
		),
	}

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

				for _, processor := range processors {
					if err := processor.Process(rssFeed); err != nil {
						logger.Println("failed to process:", err)
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
