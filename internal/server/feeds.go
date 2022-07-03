package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

type FeedsVals struct {
	MainFeeds     []feeds.Feed
	RegionalFeeds []feeds.Feed
}

func Feeds(vals FeedsVals) func(c *gin.Context) {
	type ResponseFeed struct {
		Slug           string `json:"slug"`
		Title          string `json:"title"`
		CollectionSlug string `json:"collection"`
	}

	type Response struct {
		Feeds []ResponseFeed `json:"feeds"`
	}

	return func(c *gin.Context) {
		response := Response{}

		for _, feed := range vals.MainFeeds {
			response.Feeds = append(response.Feeds, ResponseFeed{
				Slug:           feed.Slug(),
				Title:          feed.Title,
				CollectionSlug: "principali",
			})
		}

		for _, feed := range vals.RegionalFeeds {
			response.Feeds = append(response.Feeds, ResponseFeed{
				Slug:           feed.Slug(),
				Title:          feed.Title,
				CollectionSlug: "regionali",
			})
		}

		c.JSON(200, response)
	}
}
