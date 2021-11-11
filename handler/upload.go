package handler

import (
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/middleware/myjwt"
)

func Upload(c *gin.Context) {
	indexes := c.QueryArray("pos")
	claims, ok := c.Get("TokenClaims")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "无法解析token信息",
		})
		return
	}
	user := claims.(*myjwt.UserClaims).User

	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(indexes)

	var userFromDB string
	err := confs.DB.QueryRow(`SELECT user FROM datatable WHERE user=?`, user).Scan(&userFromDB)
	if err == nil {
		_, err = confs.DB.Exec(`UPDATE datatable SET data=? WHERE user=?`, buf.Bytes(), user)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "数据上传成功",
		})
	} else {
		_, err = confs.DB.Exec(`INSERT INTO datatable(user, data) VALUES (?,?)`, user, buf.Bytes())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "数据上传成功",
		})
	}
}
