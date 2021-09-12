package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kascas/httpserver/logs"
)

type Queries struct {
	Args  map[string]float64 `json:"args"`
	Scale int                `json:"scale"`
	Areas [][]float64        `json:"areas"`
}

var MyQueries *Queries

func init() {
	MyQueries = &Queries{}
	_, err := os.Stat("./.conf/query.json")
	if err != nil {
		if os.IsNotExist(err) {
			os.Create("./.conf/query.json")
		} else {
			logs.ErrorPanic(err, "获取query.json状态错误")
			return
		}
	}

	data, err := ioutil.ReadFile(`./.conf/query.json`)
	if err != nil {
		logs.ErrorPanic(err, err.Error())
		return
	}
	err = json.Unmarshal(data, MyQueries)
	if err != nil {
		logs.ErrorPanic(err, err.Error())
		return
	}
}

func Query(c *gin.Context) {
	data, err := ioutil.ReadFile(`./.conf/query.json`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	err = json.Unmarshal(data, MyQueries)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Published Areas",
		"data":   MyQueries,
	})
}