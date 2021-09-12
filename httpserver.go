package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/handler"
	"github.com/kascas/httpserver/logs"
	"github.com/kascas/httpserver/middleware/myjwt"
	"github.com/kascas/httpserver/middleware/statuslog"
	"github.com/kascas/httpserver/rappor"
	"github.com/robfig/cron"
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

	r.GET("/checklog", handler.Checklog)

	r.Use(statuslog.Statuslog)
	router := r.Group(`/user`, myjwt.JWTAuth())
	{
		router.POST(`/upload`, handler.Upload)
		router.GET(`/refresh`, func(c *gin.Context) {})
		router.GET(`/query`, handler.Query)
		router.GET(`/search`, handler.Search)
	}

	c := cron.New()
	c.AddFunc("@every 1m", rappor.StatRun)
	c.Start()

	n := confs.NetInfo
	err := r.RunTLS(n.BindIP+`:`+strconv.Itoa(int(n.Port)), `./cert.pem`, `./key.pem`)
	if err != nil {
		logs.ErrorPanic(err, `/httpserver.go`)
	}
}
