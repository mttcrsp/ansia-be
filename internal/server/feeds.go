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
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}

	type ResponseCollection struct {
		Slug  string         `json:"slug"`
		Name  string         `json:"name"`
		Feeds []ResponseFeed `json:"feeds"`
	}

	type Response struct {
		Collections []ResponseCollection `json:"collections"`
	}

	return func(c *gin.Context) {
		mainCollection := ResponseCollection{
			Slug: "principali",
			Name: "Principali",
		}
		for _, feed := range vals.MainFeeds {
			mainCollection.Feeds = append(mainCollection.Feeds, ResponseFeed{
				Slug:  feed.Slug(),
				Title: feed.Title,
			})
		}
		regionalCollection := ResponseCollection{
			Slug: "regionali",
			Name: "Regionali",
		}
		for _, feed := range vals.RegionalFeeds {
			regionalCollection.Feeds = append(regionalCollection.Feeds, ResponseFeed{
				Slug:  feed.Slug(),
				Title: feed.Title,
			})
		}
		c.JSON(200, Response{
			Collections: []ResponseCollection{mainCollection, regionalCollection},
		})
	}
}
