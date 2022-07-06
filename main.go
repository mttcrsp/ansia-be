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
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	feedsLoader := feeds.Loader{}
	rssLoader := rss.Loader{}
	store := core.Store{}

	watcherLogger := newLogger("watcher")
	watcher := core.NewWatcher(feedsLoader, rssLoader)

	extractor := articles.NewExtractor()
	extractorLogger := newLogger("extractor")
	queuedExtractor := core.NewQueuedExtractor(*extractor)

	mainFeeds, regionalFeeds, err := feedsLoader.LoadAll()
	if err != nil {
		return err
	}

	feedsHandler := server.Feeds(
		server.FeedsVals{
			MainFeeds:     mainFeeds,
			RegionalFeeds: regionalFeeds,
		},
	)

	feedBySlugHandler := server.FeedBySlug(
		server.FeedBySlugVals{
			MainFeeds:     mainFeeds,
			RegionalFeeds: regionalFeeds,
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
		watcher.Run(
			core.WatcherConfig{
				IterationBackoff: time.Minute * 5,
			},
			core.WatcherHandlers{
				OnUpdate: func(wu core.WatcherUpdate) {
					if err := store.InsertFeedItems(wu.Feed, wu.RSS); err != nil {
						watcherLogger.Println("failed to insert feed items:", err)
						return
					}

					var items []rss.Item
					for _, item := range wu.RSS.Channel.Items {
						found, err := store.ArticleExists(item.ID())
						if err != nil {
							watcherLogger.Println("failed to lookup article:", err)
							continue
						}
						if !found {
							items = append(items, item)
						}
					}
					queuedExtractor.Enqueue(items)
				},
				OnError: func(err error) {
					watcherLogger.Println(err)
				},
				OnIterationBegin: func() {
					watcherLogger.Println("iteration will begin")
				},
				OnIterationEnd: func() {
					watcherLogger.Println("iteration did end")
				},
			},
		)

		c <- "watcher did exit"
	}()

	go func() {
		queuedExtractor.Run(
			core.QueuedExtractorConfig{
				Backoff: time.Second / 2,
			},
			core.QueuedExtractorHandlers{
				OnItemExtracted: func(item core.QueuedExtractorItem) {
					extractorLogger.Printf("extracted item '%s'\n", item.Item.Link)
					if err := store.InsertArticle(item.Item, item.Article); err != nil {
						extractorLogger.Printf("failed to insert article '%s': %s\n", item.Item.Link, err)
					} else {
						extractorLogger.Printf("did insert article '%d'\n", item.Item.ID())
					}
				},
				OnError: func(err error) {
					extractorLogger.Println(err)
				},
			},
		)

		c <- "extractor did exit"
	}()

	return errors.New(<-c)
}

func newLogger(identifier string) *log.Logger {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)
	logger.SetPrefix(fmt.Sprintf("[%s] ", identifier))
	logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	return logger
}
