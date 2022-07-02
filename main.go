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

	conversionLogger := newLogger("conversion")
	toItems := func(wis []core.WatcherItem) []core.Item {
		var cis []core.Item
		for _, wi := range wis {
			ci, err := core.NewItem(wi.Item, wi.Feed)
			if err != nil {
				conversionLogger.Printf("failed to convert item '%s': %s", wi.Item.Link, err)
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

					if err := store.InsertItems(toItems(wi)); err != nil {
						storeLogger.Printf("failed to insert items: %s\n", err)
					}
					storeLogger.Printf("did insert items")

					var rssItems []rss.Item
					for _, item := range wi {
						rssItems = append(rssItems, item.Item)
					}
					queuedExtractor.Enqueue(rssItems)
				},
				OnDelete: func(wi []core.WatcherItem) {
					watcherLogger.Println("deleted", len(wi))

					if err := store.DeleteItems(toItems(wi)); err != nil {
						storeLogger.Printf("failed to delete items: %s\n", err)
					}
					storeLogger.Printf("did delete items")
				},
				OnError: func(err error) {
					watcherLogger.Println(err)
				},
				OnIterationBegin: func() {
					watcherLogger.Println("iteration did begin")
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
					storeLogger.Printf("did insert article '%s'", qei.Item.Link)
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

	feedsHandler := func(c *gin.Context) {
		type ResponseFeed struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
		}

		type ResponseCollection struct {
			Slug  string         `json:"slug"`
			Name  string         `json:"name"`
			Feeds []ResponseFeed `json:"feeds"`
		}

		type Response struct {
			Collections []ResponseCollection `json:"collections"`
		}

		mainCollection := ResponseCollection{
			Slug: "principali",
			Name: "Principali",
		}
		for _, feed := range mainFeeds {
			mainCollection.Feeds = append(mainCollection.Feeds, ResponseFeed{
				Slug:  feed.Slug(),
				Title: feed.Title,
			})
		}
		regionalCollection := ResponseCollection{
			Slug: "regionali",
			Name: "Regionali",
		}
		for _, feed := range regionalFeeds {
			regionalCollection.Feeds = append(regionalCollection.Feeds, ResponseFeed{
				Slug:  feed.Slug(),
				Title: feed.Title,
			})
		}
		c.JSON(200, Response{
			Collections: []ResponseCollection{mainCollection, regionalCollection},
		})
	}

	feedHandler := func(c *gin.Context) {
		feedSlug := c.Param("feed")

		var feed *feeds.Feed
		for _, f := range append(mainFeeds, regionalFeeds...) {
			if f.Slug() == feedSlug {
				feed = &f
			}
		}

		if feed == nil {
			c.Status(404)
			return
		}

		feedItems, err := store.GetFeed(feedSlug)
		if err != nil {
			c.Status(500)
			return
		}

		if len(feedItems) == 0 {
			c.Status(400)
			return
		}

		type Response struct {
			Items []core.FeedItem `json:"items"`
		}

		c.JSON(200, Response{Items: feedItems})
	}

	go func() {
		// gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.GET("/feeds", feedsHandler)
		r.GET("/feeds/:feed", feedHandler)
		r.Run()
		c <- "server did complete"
	}()

	fmt.Println(<-c)

	return nil
}
