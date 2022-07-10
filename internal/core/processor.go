package core

import (
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
)

type RSSProcessor interface {
	Process(feed feeds.Feed, rssFeed *rss.RSS)
}
