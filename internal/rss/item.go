package rss

import (
	"encoding/xml"
	"hash/fnv"
)

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDateRaw  string   `xml:"pubDate"`
}

func (i *Item) ID() int64 {
	h := fnv.New64()
	h.Write([]byte(i.Link))
	return int64(h.Sum64())
}
