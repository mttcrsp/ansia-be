package core

import (
	"log"
	"time"

	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/store"
)

type ArticlesProcessor struct {
	extractor articles.Extractor
	store     store.Store
	logger    log.Logger
}

func NewArticlesProcessor(extractor articles.Extractor, store store.Store, logger log.Logger) *ArticlesProcessor {
	return &ArticlesProcessor{
		extractor: extractor,
		store:     store,
		logger:    logger,
	}
}

func (p *ArticlesProcessor) Process(feed feeds.Feed, rssFeed *rss.RSS) error {
	if feed.CollectionSlug == "media" {
		return nil
	}

	for _, item := range rssFeed.Channel.Items {
		found, err := p.store.ArticleExists(item.ID())
		if err != nil {
			p.logger.Printf("failed to check availability '%d': %s", item.ID(), err)
			time.Sleep(time.Second)
			continue
		}

		if found {
			time.Sleep(time.Second)
			continue
		}

		p.logger.Println("extracting article", item.Link)
		article, err := p.extractor.Extract(item.Link)
		if err != nil {
			p.logger.Printf("failed to extract '%s': %s", item.Link, err)
			time.Sleep(time.Second)
			continue
		}

		if err = p.store.InsertArticle(item, *article); err != nil {
			p.logger.Printf("failed to insert article '%s: %s", item.Link, err)
			time.Sleep(time.Second)
			continue
		}
	}

	return nil
}
