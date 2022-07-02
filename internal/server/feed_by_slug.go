package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/core"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

type FeedBySlugVals struct {
	MainFeeds     []feeds.Feed
	RegionalFeeds []feeds.Feed
}

type FeedBySlugDeps struct {
	Store core.Store
}

func FeedBySlug(vals FeedBySlugVals, deps FeedBySlugDeps) func(c *gin.Context) {
	type Response struct {
		Items []core.FeedItem `json:"items"`
	}

	return func(c *gin.Context) {
		feedSlug := c.Param("feed")

		var feed *feeds.Feed
		for _, f := range append(vals.MainFeeds, vals.RegionalFeeds...) {
			if f.Slug() == feedSlug {
				feed = &f
			}
		}

		if feed == nil {
			c.Status(404)
			return
		}

		feedItems, err := deps.Store.GetFeed(feedSlug)
		if err != nil {
			c.Status(500)
			return
		}

		if len(feedItems) == 0 {
			c.Status(400)
			return
		}

		c.JSON(200, Response{Items: feedItems})
	}
}
