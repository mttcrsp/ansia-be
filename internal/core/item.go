package core

import (
	"hash/fnv"
	"time"

	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
)

type Item struct {
	ID          int64     `json:"item_id" db:"item_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Feed        string    `json:"feed" db:"feed"`
}

func NewItem(item rss.Item, feed feeds.Feed) (*Item, error) {
	publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDateRaw)
	if err != nil {
		return nil, err
	}

	h := fnv.New64()
	h.Write([]byte(item.Link))

	return &Item{
		ID:          int64(h.Sum64()),
		Title:       item.Title,
		Description: item.Description,
		URL:         item.Link,
		PublishedAt: publishedAt,
		Feed:        feed.Slug(),
	}, nil
}
