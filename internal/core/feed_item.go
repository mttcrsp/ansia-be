package core

import "time"

type FeedItem struct {
	ItemID      int64     `db:"item_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	URL         string    `db:"url"`
	PublishedAt time.Time `db:"published_at"`
	Feed        string    `db:"feed"`
	Keywords    string    `db:"keywords"`
	Content     string    `db:"content"`
	ImageURL    string    `db:"image_url"`
}
