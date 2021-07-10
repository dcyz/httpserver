package handler

import (
	"github.com/gin-gonic/gin"
)

// Download 选择目标文件传输给用户
func Download(c *gin.Context) {
	c.File("./file/test.png")
}
