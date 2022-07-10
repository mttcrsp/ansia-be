package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

type FeedsVals struct {
	Collections feeds.Collections
}

func Feeds(vals FeedsVals) func(c *gin.Context) {
	type ResponseFeed struct {
		Slug           string `json:"slug"`
		Title          string `json:"title"`
		Weight         int    `json:"weight"`
		CollectionSlug string `json:"collection"`
	}

	type Response struct {
		Feeds []ResponseFeed `json:"feeds"`
	}

	toResponseFeeds := func(feeds []feeds.Feed, collectionSlug string) []ResponseFeed {
		var responseFeeds []ResponseFeed
		for _, feed := range feeds {
			responseFeeds = append(responseFeeds, ResponseFeed{
				Slug:           feed.Slug(),
				Title:          feed.Title,
				Weight:         feed.Weight,
				CollectionSlug: collectionSlug,
			})
		}
		return responseFeeds
	}

	return func(c *gin.Context) {
		response := Response{}
		response.Feeds = append(response.Feeds, toResponseFeeds(vals.Collections.Main, "principali")...)
		response.Feeds = append(response.Feeds, toResponseFeeds(vals.Collections.Regional, "regionali")...)
		response.Feeds = append(response.Feeds, toResponseFeeds(vals.Collections.Media, "media")...)
		c.JSON(200, response)
	}
}
