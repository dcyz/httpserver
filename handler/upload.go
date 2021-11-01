package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/middleware/myjwt"
)

func Upload(c *gin.Context) {
	var data []byte

	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	bitLen := 0
	if int(math.Ceil(float64(bitLen)/8)) != len(data) {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "目标长度：" + strconv.FormatInt(int64(math.Ceil(float64(bitLen)/8)), 10) + " | 实际长度：" + strconv.FormatInt(int64(len(data)), 10),
		})
		return
	}
	claims, ok := c.Get("TokenClaims")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "无法解析token信息",
		})
		return
	}
	user := claims.(*myjwt.UserClaims).User

	var userFromDB string
	err = confs.DB.QueryRow(`SELECT user FROM datatable WHERE user=?`, user).Scan(&userFromDB)
	if err == nil {
		_, err = confs.DB.Exec(`UPDATE datatable SET data=? WHERE user=?`, data, user)
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
		_, err = confs.DB.Exec(`INSERT INTO datatable(user, data) VALUES (?,?)`, user, data)
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
