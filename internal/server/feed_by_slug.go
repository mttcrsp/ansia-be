package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/core"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

type FeedBySlugVals struct {
	Collections feeds.Collections
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

		found := false
		for _, f := range vals.Collections.All() {
			if f.Slug() == feedSlug {
				found = true
			}
		}

		if !found {
			c.Status(404)
			return
		}

		feedItems, err := deps.Store.GetFeed(feedSlug)
		if err != nil {
			c.Status(500)
			return
		}

		status := 200
		if len(feedItems) == 0 {
			status = 204
		}

		c.JSON(status, Response{Items: feedItems})
	}
}
