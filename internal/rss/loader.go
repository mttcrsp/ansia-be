package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Loader struct{}

func (l *Loader) Load(url string) (*RSS, error) {
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
