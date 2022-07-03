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
				OnInsert: func(items []core.WatcherItem) {
					watcherLogger.Println("found inserted", len(items))

					if err := store.InsertItems(toItems(items, watcherLogger)); err != nil {
						watcherLogger.Println("failed to insert items:", err)
						return
					}
					watcherLogger.Println("did insert items")

					var rssItems []rss.Item
					for _, item := range items {
						rssItems = append(rssItems, item.Item)
					}

					queuedExtractor.Enqueue(rssItems)
				},
				OnDelete: func(items []core.WatcherItem) {
					watcherLogger.Println("found deleted", len(items))

					if err := store.DeleteItems(toItems(items, watcherLogger)); err != nil {
						watcherLogger.Println("failed to delete items:", err)
					} else {
						watcherLogger.Println("did delete items")
					}
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

					article := core.NewArticle(item.Article, item.Item)
					if err := store.InsertArticle(article); err != nil {
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

func toItems(wis []core.WatcherItem, logger *log.Logger) []core.Item {
	var cis []core.Item
	for _, wi := range wis {
		ci, err := core.NewItem(wi.Item, wi.Feed)
		if err != nil {
			logger.Printf("failed to convert item '%s': %s\n", wi.Item.Link, err)
			continue
		}

		cis = append(cis, *ci)
	}

	return cis
}
