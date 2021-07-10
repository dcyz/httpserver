package main

import (
	"httpserver/confs"
	"httpserver/handler"
	"httpserver/logs"
	"httpserver/middleware/myjwt"
	"httpserver/middleware/tls"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init() {
	confs.WdInit()
	logs.LogInit()
	confs.NetInit()
	confs.DBInit()
	confs.CertInit()
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(tls.TLS())
	router := r.Group(`/auth`, myjwt.JWTAuth())
	{
		router.POST("/upload", handler.Upload)
		router.GET("/download", handler.Download)
	}
	r.POST(`/signup`, handler.SignUp)
	r.GET(`/signin`, handler.SignIn)

	n := confs.NetInfo
	err := r.RunTLS(n.Host+`:`+strconv.Itoa(int(n.Port)),
		`./cert.pem`, `./key.pem`)
	if err != nil {
		logs.ErrorPanic(err, `/httpserver.go`)
	}
}
