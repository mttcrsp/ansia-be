package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
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
	item, err := core.NewItem(
		rss.Item{
			Title:       "something",
			Description: "Else",
			Link:        "https://www.ansa.it/something/else",
			PubDateRaw:  "Mon, 2 Jan 2006 15:04:05 -0700",
		},
		feeds.Feed{
			Title: "Politica",
			URL:   "https://www.ansa.it/sito/notizie/politica/politica_rss.xml",
		},
	)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	store := core.Store{}
	if err := store.Insert([]core.Item{*item}); err != nil {
		return err
	}

	// newLogger := func(identifier string) log.Logger {
	// 	logger := log.Logger{}
	// 	logger.SetOutput(os.Stdout)
	// 	logger.SetPrefix(fmt.Sprintf("[%s] ", identifier))
	// 	logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// 	return logger
	// }

	// fl := feeds.Loader{}
	// rl := rss.Loader{}
	// watcherLogger := newLogger("watcher")
	// watcher := core.NewWatcher(fl, rl)

	// extractor := articles.NewExtractor()
	// extractorLogger := newLogger("extractor")
	// queuedExtractor := core.NewQueuedExtractor(*extractor)

	// c := make(chan string)

	// go func() {
	// 	watcher.Run(
	// 		core.WatcherConfig{
	// 			IterationBackoff: time.Minute,
	// 		},
	// 		core.WatcherHandlers{
	// 			OnInsert: func(wi []core.WatcherItem) {
	// 				watcherLogger.Println("inserted", len(wi))

	// 				var items []rss.Item
	// 				for _, item := range wi {
	// 					items = append(items, item.Item)
	// 				}
	// 				queuedExtractor.Enqueue(items)
	// 			},
	// 			OnDelete: func(wi []core.WatcherItem) {
	// 				watcherLogger.Println("deleted", len(wi))
	// 			},
	// 			OnError: func(err error) {
	// 				watcherLogger.Println(err)
	// 			},
	// 			OnIterationBegin: func() {
	// 				watcherLogger.Println("iteration did begin")
	// 			},
	// 			OnIterationEnd: func() {
	// 				watcherLogger.Println("iteration did end")
	// 			},
	// 		},
	// 	)
	// 	c <- "watcher did complete"
	// }()

	// go func() {
	// 	queuedExtractor.Run(
	// 		core.QueuedExtractorConfig{
	// 			Backoff: time.Second,
	// 		},
	// 		core.QueuedExtractorHandlers{
	// 			OnItemExtracted: func(qei core.QueuedExtractorItem) {
	// 				extractorLogger.Printf("extracted item '%s'\n", qei.Item.Link)
	// 			},
	// 			OnError: func(err error) {
	// 				extractorLogger.Println(err)
	// 			},
	// 		},
	// 	)
	// 	c <- "extractor did complete"
	// }()

	// fmt.Println(<-c)

	return nil
}
