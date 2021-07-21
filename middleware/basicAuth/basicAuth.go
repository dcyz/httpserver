package basicAuth

import (
	"net/http"
	"strings"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/logs"

	"github.com/gin-gonic/gin"
)

func getPasswd(user string) (string, bool) {
	var passwd string
	err := confs.DB.QueryRow(`SELECT passwd FROM authtable WHERE user=?`, user).Scan(&passwd)
	if err != nil {
		logs.ErrorLog(err, `密码获取失败`)
		return ``, false
	}
	return passwd, true
}

func BasicAuth(c *gin.Context) {
	user, passwd, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    `无Auth信息`,
		})
		c.Abort()
		return
	}
	passwdSQL, ok := getPasswd(user)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    `认证处理异常`,
		})
		c.Abort()
		return
	} else if strings.Compare(passwd, passwdSQL) != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    `认证失败`,
		})
		c.Abort()
		return
	}
}
