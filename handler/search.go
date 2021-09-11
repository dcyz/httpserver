package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Points struct {
	Lat float64
	Lng float64
}

func Search(c *gin.Context) {
	var point Points
	if c.BindJSON(&point) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `DataParseFailed`,
		})
		return
	}
	
}
