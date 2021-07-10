package handler

import (
	"httpserver/confs"
	"httpserver/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload 用户上传，路由为/auth/upload
func Upload(c *gin.Context) {
	// 从请求中获取data
	data, err := c.GetRawData()
	if err != nil {
		logs.ErrorLog(err, `Upload： 数据解析失败`)
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `数据解析异常`,
		})
		return
	}
	// 将data插入到数据库中
	_, err = confs.DB.Exec(`INSERT INTO demotable(data) VALUES (?)`, data)
	if err != nil {
		logs.ErrorLog(err, `Upload: 上传数据异常`)
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `上传失败`,
		})
		return
	}
	// 回写到用户
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    `上传成功`,
	})
}
