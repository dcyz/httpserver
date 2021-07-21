package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/middleware/myjwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// userInfo 用户信息
type userInfo struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

// auth 认证
type auth struct {
	Token string
	userInfo
}

// check 根据数据库的authtable表验证输入的user-passwd是否正确
func check(user string, passwd string) error {
	// 从DB获取user对应的passwd
	var passwdFromDB string
	err := confs.DB.QueryRow(`SELECT passwd FROM authtable WHERE user=?`, user).Scan(&passwdFromDB)
	if err != nil {
		// logs.ErrorLog(err, `signin.go -> check异常`)
		return errors.New(`用户不存在`)
	}
	// 如果passwd一致则返回nil
	if strings.Compare(passwd, passwdFromDB) == 0 {
		return nil
	} else {
		return errors.New(`用户名或密码错误`)
	}
}

// SignIn 处理用户的登录请求，并分发一个token
func SignIn(c *gin.Context) {
	var u userInfo
	// 从request的json内获取user和passwd
	if c.BindJSON(&u) == nil {
		// check函数进行验证
		err := check(u.User, u.Passwd)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		} else {
			// 如果check无错，则生成token
			generateToken(c, u)
		}
	}
}

// generateToken 根据userInfo生成Token
func generateToken(c *gin.Context, u userInfo) {
	// 新建JWT实例
	k := &myjwt.KeyStruct{
		Key: []byte(myjwt.GetSignKey()),
	}
	// 新建CustomClaims实例
	claims := myjwt.CustomClaims{
		User: u.User,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + 3600,
			Issuer:    "dcyz",
		},
	}
	// 生成新的Token
	token, err := k.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	// 将token发送给用户
	data := auth{
		Token:    token,
		userInfo: u,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    `登录成功`,
		"data":   data,
	})
}
