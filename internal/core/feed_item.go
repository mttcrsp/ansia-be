package core

import (
	"encoding/json"
	"strings"
	"time"
)

type FeedItem struct {
	ItemID      int64     `db:"item_id" json:"item_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	URL         string    `db:"url" json:"url"`
	PublishedAt time.Time `db:"published_at" json:"published_at"`
	Feed        string    `db:"feed" json:"feed"`
	Keywords    string    `db:"keywords" json:"keywords"`
	Content     string    `db:"content" json:"content"`
	ImageURL    string    `db:"image_url" json:"image_url"`
}

func (i *FeedItem) MarshalJSON() ([]byte, error) {
	type Alias FeedItem
	return json.Marshal(&struct {
		Keywords []string `json:"keywords"`
		*Alias
	}{
		Keywords: strings.Split(i.Keywords, ","),
		Alias:    (*Alias)(i),
	})
}
