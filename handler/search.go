package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/rappor"
)

func Search(c *gin.Context) {
	if rappor.Result != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "数据获取成功",
			"data":   rappor.Result,
		})
	}
}
