package articles

import (
	"html"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func (e *Extractor) Extract(articleURL string) (*Article, error) {
	article, err := e.g.ExtractFromURL(articleURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.RawHTML))
	if err != nil {
		log.Fatal(err)
	}

	node, err := doc.Find("div.news-txt p").Html()
	if err != nil {
		return nil, err
	}

	imageURL, err := url.Parse(article.TopImage)
	if err != nil {
		return nil, err
	}

	imageURL.Scheme = "https"

	rawContent := strings.ReplaceAll(node, "\n", " ")
	rawContent = strings.ReplaceAll(rawContent, "<strong>", "")
	rawContent = strings.ReplaceAll(rawContent, "</strong>", "")
	rawContent = strings.ReplaceAll(rawContent, "<br>", "\n")
	rawContent = strings.ReplaceAll(rawContent, "<br/>", "\n")
	rawContent = html.UnescapeString(rawContent)

	paragraphs := strings.Split(rawContent, "\n")
	var trimmedParagraphs []string
	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if len(trimmed) == 0 {
			continue
		}
		trimmedParagraphs = append(trimmedParagraphs, trimmed)
	}

	content := strings.Join(trimmedParagraphs, "\n")
	content = strings.TrimPrefix(content, "(ANSA) - ")
	content = strings.TrimPrefix(content, "(ANSA-AFP) - ")
	content = strings.TrimSuffix(content, "(ANSA).")
	content = strings.TrimSuffix(content, "(ANSA-AFP).")
	content = strings.TrimSpace(content)

	return &Article{
		Title:       article.Title,
		Description: article.MetaDescription,
		Keywords:    article.MetaKeywords,
		ImageURL:    imageURL.String(),
		Content:     content,
	}, nil
}
