package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/middleware/myjwt"
)

func Refresh(c *gin.Context) {
	types, ok := c.Get("TokenType")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "RefreshGetTokenTypeError",
		})
		return
	} else {
		if types.(int) != 0x1 {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "RefreshWithWrongTokenType",
			})
		}
	}

	claim, ok := c.Get("TokenClaims")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "RefreshGetClaimError",
		})
	}
	token, err := myjwt.GetRefreshToken(claim.(*myjwt.UserClaims).User)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "RefreshGetTokenError",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "RefreshSucess",
		"data":   token,
	})
}
