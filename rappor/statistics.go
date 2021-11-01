package rappor

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/logs"
)

var stat []float64
var count int
var blocks, random int
var f, p, q float64

type Areas struct {
	Areas [][]float64 `json:"areas"`
}

var Result []int
var MyAreas Areas

func setupStat() {
	count = 0
	f, p, q = 0.25, 0.25, 0.75

	raw, err := ioutil.ReadFile("./.conf/areas.json")
	if err != nil {
		logs.ErrorPanic(err, `areas.json读写失败`)
	}
	// 解序列化数据
	err = json.Unmarshal(raw, &MyAreas)
	if err != nil {
		logs.ErrorPanic(err, `areas.json解序列化失败`)
	}
	blocks, random = len(MyAreas.Areas), len(MyAreas.Areas)/2
	stat = make([]float64, blocks)
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
	var sum float64 = 0
	for i := 0; i < blocks; i++ {
		stat[i] = (stat[i] - k*(p+f*q/2-f*p/2)*n) / (1 - f) / (q - p)
		sum += stat[i]
	}
	for i := 0; i < blocks; i++ {
		stat[i] /= sum
	}
}

func dataAnalyze() {
	total := 0
	Result = make([]int, blocks)

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
		total++
	}
	compute()
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	for i := 0; i < blocks; i++ {
		Result[i] = int(math.Round(stat[i] * float64(total)))
	}
	fmt.Println(">>> 统计密度分布：", Result)
}

func StatRun() {
	setupStat()
	dataAnalyze()
}
