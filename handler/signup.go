package handler

import (
	"net/http"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/middleware/myjwt"

	"github.com/gin-gonic/gin"
)

// SignUp 用户注册
func SignUp(c *gin.Context) {
	var u userInfo
	// 从data中获取user和passwd
	if c.BindJSON(&u) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `DataParseFailed`,
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
			"msg":    `UserAlreadyExisted`,
		})
		return
	}
	// 将user:passwd插入到数据库中
	_, err = confs.DB.Exec(`INSERT INTO authtable(user, passwd) VALUES (?,?)`, u.User, u.Passwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    `SignUpFailed`,
		})
		return
	}
	autoSignIn(c, u)
}

func autoSignIn(c *gin.Context, u userInfo) {
	token, err := myjwt.GetAccessToken(u.User)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    `LoginSucess`,
		"data": auth{
			Token: token,
			User:  u.User,
		},
	})
}
