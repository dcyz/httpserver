package statuslog

import (
	"container/list"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/logs"
)

var logQueue = list.New()
var current = 0
var max = 1000

func Statuslog(c *gin.Context) {
	c.Next()
	if c.Writer.Status() == http.StatusOK ||
		(c.Writer.Status() == http.StatusUnauthorized && c.Request.RequestURI == "/user/refresh") {
		return
	}
	desc := fmt.Sprintf("[%d] | %s => (%s) %s\n",
		c.Writer.Status(),
		c.Request.RemoteAddr,
		c.Request.Method,
		c.Request.RequestURI)
	logQueue.PushBack(desc)

	if current < max {
		current++
		fp, err := os.OpenFile("./statuslog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		defer func() {
			err := fp.Close()
			if err != nil {
				logs.ErrorLog(err, err.Error())
			}
		}()
		if err != nil {
			logs.ErrorLog(err, err.Error())
		}
		fp.WriteString(desc)
	} else {
		logQueue.Remove(logQueue.Front())
		fp, err := os.OpenFile("./statuslog.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0777)
		defer func() {
			err := fp.Close()
			if err != nil {
				logs.ErrorLog(err, err.Error())
			}
		}()
		if err != nil {
			logs.ErrorLog(err, err.Error())
		}
		element := logQueue.Front()
		fp.WriteString(element.Value.(string))
		for i := 1; i < logQueue.Len(); i++ {
			element = element.Next()
			fp.WriteString(element.Value.(string))
		}
	}

}
