package core

import (
	"time"

	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/utils"
)

type WatcherUpdate struct {
	Feed feeds.Feed
	RSS  rss.RSS
}

type WatcherConfig struct {
	IterationBackoff time.Duration
}

type WatcherHandlers struct {
	OnUpdate         func(WatcherUpdate)
	OnIterationBegin func()
	OnIterationEnd   func()
	OnError          func(error)
}

type Watcher struct {
	feedsLoader feeds.Loader
	rssLoader   rss.Loader
}

func NewWatcher(feedsLoader feeds.Loader, rssLoader rss.Loader) *Watcher {
	return &Watcher{
		feedsLoader: feedsLoader,
		rssLoader:   rssLoader,
	}
}

func (w *Watcher) Run(config WatcherConfig, handlers WatcherHandlers) func() {
	if handlers.OnUpdate == nil {
		panic("must provide an update handler")
	}
	if handlers.OnError == nil {
		panic("must provide an error handler")
	}

	mainFeeds, regionalFeeds, err := w.feedsLoader.LoadAll()
	if err != nil {
		handlers.OnError(err)
		return func() {}
	}

	return utils.Loop(
		config.IterationBackoff,
		func() {
			if handlers.OnIterationBegin != nil {
				handlers.OnIterationBegin()
			}

			for _, feed := range append(mainFeeds, regionalFeeds...) {
				rss, err := w.rssLoader.Load(feed.URL)
				if err != nil {
					handlers.OnError(err)
					return
				}

				handlers.OnUpdate(
					WatcherUpdate{
						Feed: feed,
						RSS:  *rss,
					},
				)
			}

			if handlers.OnIterationEnd != nil {
				handlers.OnIterationEnd()
			}
		},
	)
}
