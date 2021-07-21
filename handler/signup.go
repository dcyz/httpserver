package handler

import (
	"net/http"

	"github.com/kascas/httpserver/confs"

	"github.com/gin-gonic/gin"
)

// SignUp 用户注册
func SignUp(c *gin.Context) {
	var u userInfo
	// 从data中获取user和passwd
	if c.BindJSON(&u) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `数据解析失败`,
		})
		return
	}
	// 临时变量
	var result string
	// 确保user不在数据库中（即用户不存在）
	err := confs.DB.QueryRow(`SELECT user FROM authtable WHERE user=?`, u.User).Scan(&result)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `用户已存在`,
		})
		return
	}
	// 将user:passwd插入到数据库中
	_, err = confs.DB.Exec(`INSERT INTO authtable(user, passwd) VALUES (?,?)`, u.User, u.Passwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `注册处理异常`,
		})
		return
	}
	// 回写用户
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    `注册成功`,
	})
}
