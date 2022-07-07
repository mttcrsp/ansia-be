package core

import (
	"testing"

	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
	"github.com/stretchr/testify/assert"
)

func TestStore_Integration(t *testing.T) {
	inputFeed := feeds.Feed{
		Title: "feed-title-1",
		URL:   "feed-url-1",
	}
	inputArticle1 := articles.Article{
		Keywords: "article-keywords-1",
		Content:  "article-content-1",
		ImageURL: "article-image_url-1",
	}
	inputArticle2 := articles.Article{
		Keywords: "article-keywords-2",
		Content:  "article-content-2",
		ImageURL: "article-image_url-2",
	}
	inputRSSItem1 := rss.Item{
		Title:       "item-title-1",
		Description: "item-description-1",
		Link:        "item-link-1",
		PubDateRaw:  "Mon, 2 Jan 2006 15:04:05 -0700",
	}
	inputRSSItem2 := rss.Item{
		Title:       "item-title-2",
		Description: "item-description-2",
		Link:        "item-link-2",
		PubDateRaw:  "Mon, 2 Jan 2006 15:04:06 -0700",
	}
	inputRSS1 := rss.RSS{
		Channel: rss.Channel{
			Items: []rss.Item{inputRSSItem1},
		},
	}
	inputRSS2 := rss.RSS{
		Channel: rss.Channel{
			Items: []rss.Item{inputRSSItem2},
		},
	}

	s := &Store{}

	err := s.InsertFeedItems(inputFeed, inputRSS1)
	assert.Nil(t, err)

	result, err := s.GetFeed(inputFeed.Slug())
	assert.Nil(t, err)
	assert.Equal(t, []FeedItem{}, result)

	found, err := s.ArticleExists(inputRSSItem1.ID())
	assert.Nil(t, err)
	assert.False(t, found)

	err = s.InsertArticle(inputRSSItem1, inputArticle1)
	assert.Nil(t, err)

	found, err = s.ArticleExists(inputRSSItem1.ID())
	assert.Nil(t, err)
	assert.True(t, found)

	result, err = s.GetFeed(inputFeed.Slug())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(-8237185718998111525), result[0].ItemID)
	assert.Equal(t, "item-title-1", result[0].Title)
	assert.Equal(t, "item-description-1", result[0].Description)
	assert.Equal(t, "item-link-1", result[0].URL)
	assert.Equal(t, "feed-title-1", result[0].Feed)
	assert.Equal(t, "article-keywords-1", result[0].Keywords)
	assert.Equal(t, "article-content-1", result[0].Content)
	assert.Equal(t, "article-image_url-1", result[0].ImageURL)

	err = s.InsertFeedItems(inputFeed, inputRSS2)
	assert.Nil(t, err)

	result, err = s.GetFeed(inputFeed.Slug())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	err = s.InsertArticle(inputRSSItem2, inputArticle2)
	assert.Nil(t, err)

	result, err = s.GetFeed(inputFeed.Slug())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "item-title-2", result[0].Title)

	err = s.InsertFeedItems(inputFeed, inputRSS2)
	assert.Nil(t, err)
}
