package feeds

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type FeedRSSLoader struct{}

func (l *FeedRSSLoader) Feed(url string) (*FeedRSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := FeedRSS{}
	if err = xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
