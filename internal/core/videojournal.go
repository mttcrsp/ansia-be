package core

import (
	"log"
	"strings"

	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/store"
	"github.com/mttcrsp/ansiabe/internal/videojournal"
)

type VideojournalProcessor struct {
	extractor videojournal.Extractor
	store     store.Store
	logger    log.Logger
}

func NewVideojournalProcessor(extractor videojournal.Extractor, store store.Store, logger log.Logger) *VideojournalProcessor {
	return &VideojournalProcessor{
		extractor: extractor,
		store:     store,
		logger:    logger,
	}
}

func (p *VideojournalProcessor) Process(feed feeds.Feed, rssFeed *rss.RSS) {
	if feed.Slug() != "video" {
		return
	}

	for _, item := range rssFeed.Channel.Items {
		if !strings.Contains(item.Link, "videogiornale") {
			continue
		}

		found, err := p.store.VideojournalExists(item.ID())
		if err != nil {
			p.logger.Printf("failed to check availability '%d': %s", item.ID(), err)
			continue
		}

		if found {
			continue
		}

		url, err := p.extractor.Extract(item.Link)
		if err != nil {
			p.logger.Printf("failed to extract '%s': %s", item.Link, err)
			continue
		}

		if err = p.store.InsertVideojournal(item, url); err != nil {
			p.logger.Printf("failed to insert '%s: %s", item.Link, err)
			continue
		}
	}
}
