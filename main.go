package main

import (
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
	newLogger := func(identifier string) *log.Logger {
		logger := &log.Logger{}
		logger.SetOutput(os.Stdout)
		logger.SetPrefix(fmt.Sprintf("[%s] ", identifier))
		logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
		return logger
	}

	toItems := func(wis []core.WatcherItem, logger *log.Logger) []core.Item {
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

	fl := feeds.Loader{}
	rl := rss.Loader{}
	watcherLogger := newLogger("watcher")
	watcher := core.NewWatcher(fl, rl)

	extractor := articles.NewExtractor()
	extractorLogger := newLogger("extractor")
	queuedExtractor := core.NewQueuedExtractor(*extractor)

	storeLogger := newLogger("store")
	store := core.Store{}

	c := make(chan string)

	go func() {
		watcher.Run(
			core.WatcherConfig{
				IterationBackoff: time.Minute,
			},
			core.WatcherHandlers{
				OnInsert: func(wi []core.WatcherItem) {
					watcherLogger.Println("inserted", len(wi))

					if err := store.InsertItems(toItems(wi, watcherLogger)); err != nil {
						storeLogger.Println("failed to insert items:", err)
					}
					storeLogger.Println("did insert items")

					var rssItems []rss.Item
					for _, item := range wi {
						rssItems = append(rssItems, item.Item)
					}
					queuedExtractor.Enqueue(rssItems)
				},
				OnDelete: func(wi []core.WatcherItem) {
					watcherLogger.Println("deleted", len(wi))

					if err := store.DeleteItems(toItems(wi, watcherLogger)); err != nil {
						storeLogger.Println("failed to delete items:", err)
					}
					storeLogger.Println("did delete items")
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
		c <- "watcher did complete"
	}()

	go func() {
		queuedExtractor.Run(
			core.QueuedExtractorConfig{
				Backoff: time.Second / 2,
			},
			core.QueuedExtractorHandlers{
				OnItemExtracted: func(qei core.QueuedExtractorItem) {
					extractorLogger.Printf("extracted item '%s'\n", qei.Item.Link)

					article := core.NewArticle(qei.Article, qei.Item)
					if err := store.InsertArticle(article); err != nil {
						storeLogger.Printf("failed to insert article '%d': %s\n", article.ItemID, err)
					}
					storeLogger.Printf("did insert article '%s'\n", qei.Item.Link)
				},
				OnError: func(err error) {
					extractorLogger.Println(err)
				},
			},
		)
		c <- "extractor did complete"
	}()

	feedsLogger := newLogger("server")

	mainFeeds, err := fl.LoadMain()
	if err != nil {
		feedsLogger.Println("failed to load main feeds:", err)
		return err
	}

	regionalFeeds, err := fl.LoadRegional()
	if err != nil {
		feedsLogger.Println("failed to load main feeds:", err)
		return err
	}

	feedHandler := server.FeedBySlug(
		server.FeedBySlugVals{
			MainFeeds:     mainFeeds,
			RegionalFeeds: regionalFeeds,
		},
		server.FeedBySlugDeps{
			Store: store,
		},
	)

	feedsHandler := server.Feeds(
		server.FeedsVals{
			MainFeeds:     mainFeeds,
			RegionalFeeds: regionalFeeds,
		},
	)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.GET("/feeds", feedsHandler)
		r.GET("/feeds/:feed", feedHandler)
		r.Run()
		c <- "server did complete"
	}()

	fmt.Println(<-c)

	return nil
}
