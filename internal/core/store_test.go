package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_Integration(t *testing.T) {
	time := time.Now()

	inputArticle := Article{
		ItemID:   1,
		Keywords: "keywords",
		Content:  "content",
		ImageURL: "image_url",
	}
	inputItem1 := Item{
		ID:          1,
		Title:       "title-1",
		Description: "description-1",
		URL:         "url-1",
		PublishedAt: time,
		Feed:        "feed-1",
	}
	inputItem2 := Item{
		ID:          2,
		Title:       "title-2",
		Description: "description-2",
		URL:         "url-2",
		PublishedAt: time,
		Feed:        "feed-2",
	}
	inputItems := []Item{inputItem1, inputItem2}

	s := &Store{}
	err := s.InsertItems(inputItems)
	assert.Nil(t, err)

	items, err := s.GetItems()
	assert.Nil(t, err)
	assert.Equal(t, len(inputItems), len(items))
	assert.Equal(t, inputItems[0].ID, items[0].ID)
	assert.Equal(t, inputItems[0].Title, items[0].Title)
	assert.Equal(t, inputItems[0].Description, items[0].Description)
	assert.Equal(t, inputItems[0].URL, items[0].URL)
	assert.Equal(t, inputItems[0].Feed, items[0].Feed)

	err = s.InsertArticle(inputArticle)
	assert.Nil(t, err)

	articles, err := s.GetArticles()
	assert.Equal(t, 1, len(articles))

	feedItems, err := s.GetFeed("feed-1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(feedItems))
	assert.Equal(t, inputItems[0].ID, feedItems[0].ItemID)
	assert.Equal(t, inputItems[0].Title, feedItems[0].Title)
	assert.Equal(t, inputItems[0].Description, feedItems[0].Description)
	assert.Equal(t, inputItems[0].URL, feedItems[0].URL)
	assert.Equal(t, inputItems[0].Feed, feedItems[0].Feed)
	assert.Equal(t, "keywords", feedItems[0].Keywords)
	assert.Equal(t, "content", feedItems[0].Content)
	assert.Equal(t, "image_url", feedItems[0].ImageURL)

	err = s.DeleteItems([]Item{inputItems[0]})
	assert.Nil(t, err)

	items, err = s.GetItems()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(items))

	articles, err = s.GetArticles()
	assert.Equal(t, 0, len(articles))

	err = s.DeleteItems([]Item{inputItems[1]})
	assert.Nil(t, err)

	items, err = s.GetItems()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(items))
}
