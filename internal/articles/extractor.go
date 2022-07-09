package articles

import (
	"strings"

	goose "github.com/advancedlogic/GoOse"
	"github.com/sundy-li/html2article"
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

	extracted, err := html2article.NewFromHtml(article.RawHTML)
	if err != nil {
		return nil, err
	}

	alternateArticle, err := extracted.ToArticle()
	if err != nil {
		return nil, err
	}

	content := alternateArticle.Content
	content = strings.TrimPrefix(content, "(ANSA) - ")
	content = strings.TrimSuffix(content, " (ANSA).")
	content = strings.TrimPrefix(content, "(ANSA-AFP) - ")
	content = strings.TrimSuffix(content, " (ANSA-AFP).")

	return &Article{
		Title:       article.Title,
		Description: article.MetaDescription,
		Keywords:    article.MetaKeywords,
		ImageURL:    article.TopImage,
		Content:     content,
	}, nil
}
