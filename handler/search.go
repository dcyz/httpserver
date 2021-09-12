package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AreaStat struct {
	Lat   float64
	Lng   float64
	Width float64
	Count int
}

var MyStat []AreaStat

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "数据获取成功",
		"data":   MyStat,
	})
}
