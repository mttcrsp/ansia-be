package rss

import (
	"encoding/xml"
)

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDateRaw  string   `xml:"pubDate"`
}
