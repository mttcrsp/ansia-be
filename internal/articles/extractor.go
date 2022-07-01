package articles

import (
	"strings"

	goose "github.com/advancedlogic/GoOse"
)

type Extractor struct {
	g goose.Goose
}

func NewExtractor() *Extractor {
	return &Extractor{
		g: goose.New(),
	}
}

func (e *Extractor) Extract(url string) (*Article, error) {
	article, err := e.g.ExtractFromURL(url)
	if err != nil {
		return nil, err
	}

	return &Article{
		Title:       article.Title,
		Description: article.MetaDescription,
		Keywords:    strings.Split(article.MetaKeywords, ","),
		Content:     article.CleanedText,
		ImageURL:    article.TopImage,
	}, nil
}
