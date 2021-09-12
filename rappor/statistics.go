package rappor

import (
	"fmt"
	"log"
	"math"

	"github.com/kascas/httpserver/confs"
	"github.com/kascas/httpserver/handler"
)

var stat []float64
var count int
var bitLen, block, scale int
var f, p, q float64

func setupStat() {
	count = 0
	queries := handler.MyQueries
	f, p, q = queries.Args["f"], queries.Args["p"], queries.Args["q"]
	block, scale = len(queries.Areas), queries.Scale
	bitLen = block * scale
	stat = make([]float64, bitLen)
}

func testBit(b []byte, pos int) bool {
	outer, inner := pos>>3, pos&0x7
	mask := (byte)(1 << (7 - inner))
	return (b[outer] & mask) != 0
}

func statPerData(data []byte) {
	for i := 0; i < bitLen; i++ {
		if testBit(data, i) {
			stat[i]++
		}
	}
	count++
}

func compute() {
	n := float64(count)
	var sum float64 = 0
	for i := 0; i < bitLen; i++ {
		stat[i] = 1.0 / (1 - f) * ((stat[i]-p*n)/(q-p) - (f*n)/2)
		sum += stat[i]
	}
	for i := 0; i < bitLen; i++ {
		stat[i] /= sum
	}
}

func dataAnalyze() []int {
	total := 0
	result := make([]int, block)

	rows, err := confs.DB.Query("SELECT * FROM datatable")
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var data []byte
		var user string
		rows.Scan(&user, &data)
		statPerData(data)
		total++
	}
	compute()
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	for i := 0; i < block; i++ {
		for j := 1; j < scale; j++ {
			stat[i] += stat[i+block*j]
		}
		//actual[i] = actual[i] / float64(total)
		result[i] = int(math.Round(stat[i] * float64(total)))
	}
	fmt.Println("\n>>> 统计密度分布：", result, count, total)
	fmt.Println()
	return result
}

func storeStatData(result []int) {
	handler.MyStat = make([]handler.AreaStat, block)
	for i := 0; i < block; i++ {
		handler.MyStat[i].Lng = handler.MyQueries.Areas[i][0]
		handler.MyStat[i].Lat = handler.MyQueries.Areas[i][1]
		handler.MyStat[i].Width = handler.MyQueries.Areas[i][2]
		handler.MyStat[i].Count = result[i]
	}
}

func StatRun() {
	setupStat()
	result := dataAnalyze()
	storeStatData(result)
}
