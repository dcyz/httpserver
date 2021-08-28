package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	claim, ok := c.Get("TokenClaims")
	if ok {
		c.JSON(http.StatusOK, claim)
	}
}
