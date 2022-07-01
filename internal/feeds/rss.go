package feeds

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}
