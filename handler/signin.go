package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/middleware/myjwt"

	"github.com/gin-gonic/gin"
)

// userInfo 用户信息
type userInfo struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

// auth 认证
type auth struct {
	AccessToken  string
	RefreshToken string
}

// check 根据数据库的authtable表验证输入的user-passwd是否正确
func check(user string, passwd string) error {
	// 从DB获取user对应的passwd
	var passwdFromDB string
	err := confs.DB.QueryRow(`SELECT passwd FROM authtable WHERE user=?`, user).Scan(&passwdFromDB)
	if err != nil {
		// logs.ErrorLog(err, `signin.go -> check异常`)
		return errors.New(`UserNotFound`)
	}
	// 如果passwd一致则返回nil
	if strings.Compare(passwd, passwdFromDB) == 0 {
		return nil
	} else {
		return errors.New(`UserInfoWrong`)
	}
}

// SignIn 处理用户的登录请求，并分发一个token
func SignIn(c *gin.Context) {
	var u userInfo
	// 从request的json内获取user和passwd
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
	}
	// check函数进行验证
	err = check(u.User, u.Passwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	} else {
		// 如果check无错，则生成token
		accessToken, err := myjwt.GetAccessToken(u.User)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		refreshToken, err := myjwt.GetRefreshToken(u.User)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    `LoginSuccess`,
			"data": auth{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		})
	}
}
