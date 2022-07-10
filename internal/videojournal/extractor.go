package videojournal

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type Extractor struct{}

func (e *Extractor) Extract(itemURL string) (string, error) {
	resp, err := http.Get(itemURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var src string
	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if src != "" {
			return
		}

		if node.DataAtom.String() == "source" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					src = attr.Val
					return
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(doc)

	videoURL, err := url.Parse(src)
	if err != nil {
		return "", err
	}

	videoURL.Scheme = "https"
	return videoURL.String(), nil
}
