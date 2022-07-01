package core

import (
	"sync"
	"time"

	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/mttcrsp/ansiabe/internal/utils"
)

type QueuedExtractorItem struct {
	Item    rss.Item
	Article articles.Article
}

type QueuedExtractorConfig struct {
	Backoff time.Duration
}

type QueuedExtractorHandlers struct {
	OnError         func(error)
	OnItemExtracted func(QueuedExtractorItem)
}

type QueuedExtractor struct {
	queue     []rss.Item
	queueMu   sync.Mutex
	extractor articles.Extractor
}

func NewQueuedExtractor(extractor articles.Extractor) *QueuedExtractor {
	return &QueuedExtractor{
		extractor: extractor,
	}
}

func (e *QueuedExtractor) Run(config QueuedExtractorConfig, handlers QueuedExtractorHandlers) func() {
	return utils.Loop(
		config.Backoff,
		func() {
			e.queueMu.Lock()
			var item *rss.Item
			if len(e.queue) > 0 {
				item = &e.queue[0]
				e.queue = e.queue[1:]
			}
			e.queueMu.Unlock()

			if item == nil {
				return
			}

			article, err := e.extractor.Extract(item.Link)
			if err != nil {
				handlers.OnError(err)
				return
			}

			handlers.OnItemExtracted(
				QueuedExtractorItem{
					Item:    *item,
					Article: *article,
				},
			)
		},
	)
}

func (e *QueuedExtractor) Enqueue(items []rss.Item) {
	e.queueMu.Lock()
	e.queue = append(e.queue, items...)
	e.queueMu.Unlock()
}
