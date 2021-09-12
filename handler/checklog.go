package handler

import "github.com/gin-gonic/gin"

func Checklog(c *gin.Context) {
	c.File("statuslog.log")
}
