package core

import "github.com/mttcrsp/ansiabe/internal/articles"

type Article struct {
	Keywords string `db:"keywords"`
	Content  string `db:"content"`
	ImageURL string `db:"image_url"`
}

func NewArticle(article articles.Article) Article {
	return Article{
		Keywords: article.Keywords,
		Content:  article.Content,
		ImageURL: article.ImageURL,
	}
}