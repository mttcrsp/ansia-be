package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/feeds"
)

type FeedsVals struct {
	Collections feeds.Collections
}

func Feeds(vals FeedsVals) func(c *gin.Context) {
	type Response struct {
		Feeds []feeds.Feed `json:"feeds"`
	}

	return func(c *gin.Context) {
		c.JSON(200, Response{Feeds: vals.Collections.All()})
	}
}
