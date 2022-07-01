package feeds

import "encoding/xml"

type FeedRSS struct {
	XMLName xml.Name    `xml:"rss"`
	Channel FeedChannel `xml:"channel"`
}

type FeedChannel struct {
	XMLName xml.Name   `xml:"channel"`
	Item    []FeedItem `xml:"item"`
}
