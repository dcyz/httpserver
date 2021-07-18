package main

import (
	"httpserver/confs"
	"httpserver/handler"
	"httpserver/logs"
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
	// TODO 在此处进行其他路由
	//router := r.Group(`/`, myjwt.JWTAuth()){}
	r.POST(`/signup`, handler.SignUp)
	r.POST(`/signin`, handler.SignIn)

	n := confs.NetInfo
	err := r.RunTLS(n.BindIP+`:`+strconv.Itoa(int(n.Port)), `./cert.pem`, `./key.pem`)
	if err != nil {
		logs.ErrorPanic(err, `/httpserver.go`)
	}
}
