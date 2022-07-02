package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_Integration(t *testing.T) {
	time := time.Now()
	inputItems := []Item{
		{
			ID:          1,
			Title:       "title-1",
			Description: "description-1",
			URL:         "url-1",
			PublishedAt: time,
			Feed:        "feed-1",
		},
		{
			ID:          2,
			Title:       "title-2",
			Description: "description-2",
			URL:         "url-2",
			PublishedAt: time,
			Feed:        "feed-2",
		},
	}

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

	err = s.DeleteItems([]Item{inputItems[0]})
	assert.Nil(t, err)

	items, err = s.GetItems()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(items))

	err = s.DeleteItems([]Item{inputItems[1]})
	assert.Nil(t, err)

	items, err = s.GetItems()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(items))
}
