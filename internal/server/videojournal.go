package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mttcrsp/ansiabe/internal/store"
)

type VideojournalDeps struct {
	Store store.Store
}

func Videojournal(deps VideojournalDeps) func(c *gin.Context) {
	type Response struct {
		Videojournals []store.Videojournal `json:"videojournals"`
	}

	return func(c *gin.Context) {
		videojournals, err := deps.Store.GetVideojournals()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, Response{Videojournals: videojournals})
	}
}
