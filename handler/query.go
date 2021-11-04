package handler

import "github.com/gin-gonic/gin"

func Query(c *gin.Context) {
	c.File("./.conf/data.json")
}
