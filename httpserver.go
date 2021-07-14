package main

import (
	"httpserver/confs"
	"httpserver/handler"
	"httpserver/logs"
	"httpserver/middleware/myjwt"
	"strconv"

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
	router := r.Group(`/`, myjwt.JWTAuth())
	{
		router.POST("/upload", handler.Upload)
		router.GET("/download", handler.Download)
	}
	r.POST(`/signup`, handler.SignUp)
	r.GET(`/signin`, handler.SignIn)

	n := confs.NetInfo
	err := r.RunTLS(n.IP+`:`+strconv.Itoa(int(n.Port)), `./cert.pem`, `./key.pem`)
	if err != nil {
		logs.ErrorPanic(err, `/httpserver.go`)
	}
}
