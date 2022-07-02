package core

import (
	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/rss"
)

type Article struct {
	ItemID   int64  `db:"item_id"`
	Keywords string `db:"keywords"`
	Content  string `db:"content"`
	ImageURL string `db:"image_url"`
}

func NewArticle(article articles.Article, item rss.Item) Article {
	return Article{
		ItemID:   item.ID(),
		Keywords: article.Keywords,
		Content:  article.Content,
		ImageURL: article.ImageURL,
	}
}
