package feeds

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type FeedItem struct {
	XMLName     xml.Name `xml:"item" json:"-"`
	Title       string   `xml:"title" json:"title"`
	Description string   `xml:"description" json:"description"`
	Link        string   `xml:"link" json:"link"`
	PubDateRaw  string   `xml:"pubDate" json:"pubDate"`
}

func (i *FeedItem) MarshalJSON() ([]byte, error) {
	publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", i.PubDateRaw)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Title       string `json:"title"`
		Description string `json:"headline"`
		Link        string `json:"url"`
		PubDate     time.Time `json:"published_at"`
	}{
		Title:       i.Title,
		Description: i.Description,
		Link:        i.Link,
		PubDate:     publishedAt,
	})
}
