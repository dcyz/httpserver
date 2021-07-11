package tls

import (
	"httpserver/confs"
	"httpserver/logs"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func TLS() gin.HandlerFunc {
	return func(c *gin.Context) {
		n := confs.NetInfo
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     n.Host + `:` + strconv.Itoa(int(n.Port)),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			logs.ErrorPanic(err, `TLS 异常`)
		}
		c.Next()
	}
}
