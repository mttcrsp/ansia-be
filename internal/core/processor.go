package core

import "github.com/mttcrsp/ansiabe/internal/rss"

type RSSProcessor interface {
	Process(rssFeed *rss.RSS) error
}
