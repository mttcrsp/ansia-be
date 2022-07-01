package core

import (
	"sync"
	"time"

	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
)

type WatcherItem struct {
	Item rss.Item
	Feed feeds.Feed
}

type WatcherConfig struct {
	IterationBackoff time.Duration
}

type WatcherHandlers struct {
	OnInsert         func([]WatcherItem)
	OnDelete         func([]WatcherItem)
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

func (u *Watcher) Run(config WatcherConfig, handlers WatcherHandlers) func() {
	if handlers.OnError == nil {
		panic("must provide an error handler")
	}

	cancelled := false
	var mu sync.Mutex

	cancel := func() {
		mu.Lock()
		cancelled = true
		mu.Unlock()
	}

	isCancelled := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return cancelled
	}

	mainFeeds, err := u.feedsLoader.LoadMain()
	if err != nil {
		handlers.OnError(err)
		return cancel
	}

	regionalFeeds, err := u.feedsLoader.LoadRegional()
	if err != nil {
		handlers.OnError(err)
		return cancel
	}

	oldItems := map[string]WatcherItem{}
	iterate := func() {
		if handlers.OnIterationBegin != nil {
			handlers.OnIterationBegin()
		}

		newItems := map[string]WatcherItem{}
		for _, feed := range append(mainFeeds, regionalFeeds...) {
			rss, err := u.rssLoader.Load(feed.URL)
			if err != nil {
				handlers.OnError(err)
				return
			}

			for _, item := range rss.Channel.Items {
				newItems[item.Link] = WatcherItem{
					Feed: feed,
					Item: item,
				}
			}
		}

		if handlers.OnDelete != nil {
			var deletedItems []WatcherItem
			for link, item := range oldItems {
				if _, exists := newItems[link]; !exists {
					deletedItems = append(deletedItems, item)
				}
			}
			if deletedItems != nil {
				handlers.OnDelete(deletedItems)
			}
		}

		if handlers.OnInsert != nil {
			var insertedItems []WatcherItem
			for link, item := range newItems {
				if _, exists := oldItems[link]; !exists {
					insertedItems = append(insertedItems, item)
				}
			}
			if insertedItems != nil {
				handlers.OnInsert(insertedItems)
			}
		}

		if handlers.OnIterationEnd != nil {
			handlers.OnIterationEnd()
		}

		oldItems = newItems
	}

	for !isCancelled() {
		iterate()
		time.Sleep(config.IterationBackoff)
	}

	return cancel
}
