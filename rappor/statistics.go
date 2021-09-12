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

func SetupStat() {
	count = 0
	args := handler.MyQueries
	f, p, q = args.Args["f"], args.Args["p"], args.Args["q"]
	block, scale = len(args.Areas), args.Scale
	bitLen = block * scale
	stat = make([]float64, bitLen)
}

func testBit(b []byte, pos int) bool {
	outer, inner := pos>>3, pos&0x7
	mask := (byte)(1 << (7 - inner))
	return (b[outer] & mask) != 0
}

func StatPerData(data []byte) {
	for i := 0; i < bitLen; i++ {
		if testBit(data, i) {
			stat[i]++
		}
	}
	count++
}

func Compute() {
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

func GetResultOfLoc(loc int) float64 {
	return stat[loc]
}

func GetResult() []float64 {
	return stat
}

func DataAnalyze() {
	total := 0
	stat := make([]float64, block)

	var result []float64

	rows, err := confs.DB.Query("SELECT * FROM datatable")
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var data []byte
		var user string
		rows.Scan(&user, &data)
		StatPerData(data)
		total++
	}
	Compute()
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	result = GetResult()
	for i := 0; i < block; i++ {
		for j := 0; j < scale; j++ {
			stat[i] += result[i+block*j]
		}
		//actual[i] = actual[i] / float64(total)
		stat[i] = math.Round(stat[i] * float64(total))
	}
	fmt.Println("\n>>> 统计密度分布：", stat)
	fmt.Println()
}
