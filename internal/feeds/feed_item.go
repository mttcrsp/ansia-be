package feeds

import "encoding/xml"

type FeedItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDateRaw  string   `xml:"pubDate"`
}
