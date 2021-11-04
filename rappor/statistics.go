package rappor

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/logs"
)

var stat []float64
var count int
var blocks, random int
var f, p, q, k float64

type Areas struct {
	Areas [][]float64        `json:"areas"`
	Args  map[string]float64 `json:"args"`
}

var Result []int
var tmp []int
var MyAreas Areas

func setupStat() {
	count = 0

	raw, err := ioutil.ReadFile("./.conf/data.json")
	if err != nil {
		logs.ErrorPanic(err, `data.json读写失败`)
	}
	// 解序列化数据
	err = json.Unmarshal(raw, &MyAreas)
	if err != nil {
		logs.ErrorPanic(err, `data.json解序列化失败`)
	}

	f, p, q, k = MyAreas.Args["f"], MyAreas.Args["p"], MyAreas.Args["q"], MyAreas.Args["k"]
	blocks = len(MyAreas.Areas)
	random = int(float64(blocks) * k)
	stat = make([]float64, blocks)
	fmt.Println("")
	fmt.Println(">>> 当前时间:", time.Now())
	fmt.Println(">>> 统计参数:", f, p, q, k, blocks, random)
}

func statPerData(index []int) {
	for i := 0; i < len(index); i++ {
		stat[index[i]]++
	}
	count++
}

func compute() {
	n := float64(count)
	k := float64(random) / float64(blocks)
	// var sum float64 = 0
	for i := 0; i < blocks; i++ {
		stat[i] = (stat[i] - k*(p+f*q/2-f*p/2)*n) / (1 - f) / (q - p)
	}
}

func dataAnalyze() {
	Result = make([]int, blocks)
	tmp = make([]int, blocks)

	rows, err := confs.DB.Query("SELECT * FROM datatable")
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var data []byte
		var user string
		var index []int
		rows.Scan(&user, &data)

		gob.NewDecoder(bytes.NewReader(data)).Decode(&index)
		statPerData(index)
	}
	compute()
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	for i := 0; i < blocks; i++ {
		tmp[i] = int(math.Round(stat[i]))
	}
	Result = tmp
	fmt.Println(">>> 统计密度分布：", Result)
	fmt.Println("")
}

func StatTask() {
	setupStat()
	dataAnalyze()
}
