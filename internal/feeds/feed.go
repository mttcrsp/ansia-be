package feeds

import "github.com/gosimple/slug"

type Feed struct {
	Title          string `json:"title"`
	URL            string `json:"url"`
	Weight         int    `json:"weight"`
	CollectionSlug string `json:"collection"`
}

func (f *Feed) Slug() string {
	return slug.Make(f.Title)
}
