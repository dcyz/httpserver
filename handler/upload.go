package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/middleware/myjwt"
)

func Upload(c *gin.Context) {
	claim, ok := c.Get("TokenClaims")
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    claim.(*myjwt.UserClaims).User,
		})
	}
}
