package main

import (
	"fmt"
	"log"
	"time"

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
	var cancel func()

	fl := feeds.Loader{}
	rl := rss.Loader{}
	watcher := core.NewWatcher(fl, rl)
	cancel = watcher.Run(
		core.WatcherConfig{
			IterationBackoff: time.Minute,
		},
		core.WatcherHandlers{
			OnIterationBegin: func() {
				fmt.Println("iteration began")
			},
			OnIterationEnd: func() {
				fmt.Println("iteration ended")
			},
			OnInsert: func(wi []core.WatcherItem) {
				fmt.Println("inserted", len(wi))
			},
			OnDelete: func(wi []core.WatcherItem) {
				fmt.Println("deleted", len(wi))
			},
			OnError: func(err error) {
				if cancel != nil {
					cancel()
				}
				fmt.Println(err)
			},
		},
	)

	return nil
}
