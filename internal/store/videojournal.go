package store

type Videojournal struct {
	VideojournalID int64  `db:"item_id" json:"videojournal_id"`
	Title          string `db:"title" json:"title"`
	VideoURL       string `db:"video_url" json:"video_url"`
	PublishedAt    string `db:"published_at" json:"published_at"`
}
