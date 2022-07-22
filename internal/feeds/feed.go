package feeds

import (
	"encoding/json"

	"github.com/gosimple/slug"
)

type Feed struct {
	Title          string `json:"title"`
	URL            string `json:"url"`
	Emoji          string `json:"emoji"`
	Weight         int    `json:"weight"`
	CollectionSlug string `json:"collection"`
}

func (f *Feed) Slug() string {
	return slug.Make(f.Title)
}

func (f *Feed) MarshalJSON() ([]byte, error) {
	type Alias Feed
	return json.Marshal(&struct {
		Slug string `json:"slug"`
		*Alias
	}{
		Slug:  f.Slug(),
		Alias: (*Alias)(f),
	})
}
