package feeds

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type RSSLoader struct{}

func (l *RSSLoader) Load(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := RSS{}
	if err = xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
