package main

import (
	"strconv"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/handler"
	"github.com/kascas/httpserver/logs"
	"github.com/kascas/httpserver/middleware/myjwt"

	"github.com/gin-gonic/gin"
)

func init() {
	confs.WdInit()
	confs.NetInit()
	confs.DBInit()
	confs.CertInit()
	logs.LogInit()
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// TODO 在此处进行其他路由
	r.POST(`/signup`, handler.SignUp)
	r.POST(`/signin`, handler.SignIn)
	router := r.Group(`/user`, myjwt.JWTAuth())
	{
		router.POST(`/upload`, handler.Upload)
		router.GET(`/refresh`, func(c *gin.Context) {})
		router.GET(`/query`, handler.Query)
	}

	n := confs.NetInfo
	err := r.RunTLS(n.BindIP+`:`+strconv.Itoa(int(n.Port)), `./cert.pem`, `./key.pem`)
	if err != nil {
		logs.ErrorPanic(err, `/httpserver.go`)
	}
}
