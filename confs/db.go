package confs

import (
	"database/sql"
	"encoding/json"
	"httpserver/logs"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type DBConf struct {
	IP     string
	Port   uint32
	User   string
	Passwd string
	DBName string

	MaxOpenConns int
	MaxIdleConns int
}

// DBInfo 全局变量，数据库属性
var DBInfo *DBConf
var DB *sql.DB

// DBInit 通过db.json初始化数据库属性
func DBInit() {
	DBInfo = &DBConf{
		IP:     "0.0.0.0",
		Port:   3306,
		User:   "",
		Passwd: "",
		DBName: "",

		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	Reload()
	Connect()
	CreateTables()
}

// TODO : User and Passwd should be input by administrator
// Reload 从db.json读取数据库相关设置
func Reload() {
	data, err := ioutil.ReadFile("./.conf/db.json")
	if err != nil {
		logs.ErrorPanic(err, `数据库配置文件读写失败`)
	}
	// 解序列化数据
	err = json.Unmarshal(data, DBInfo)
	if err != nil {
		logs.ErrorPanic(err, `数据库配置文件解序列化失败`)
	}
}

func Connect() {
	var err error
	s := DBInfo
	// 建立dataSrcName
	dataSrcName := s.User + ":" + s.Passwd + "@tcp(" + s.IP + ":" + strconv.Itoa(int(s.Port)) + ")/" + s.DBName
	// 建立一个连接
	DB, err = sql.Open("mysql", dataSrcName)
	if err != nil {
		logs.ErrorPanic(err, `数据库连接失败`)
	}
	// ping测试sqlserver
	err = DB.Ping()
	if err != nil {
		logs.ErrorPanic(err, `数据库ping测试失败`)
	}
	// 设置MaxOpenConns和SetMaxIdleConns
	DB.SetMaxOpenConns(s.MaxOpenConns)
	DB.SetMaxIdleConns(s.MaxIdleConns)
}

func CreateTables() {
	_, err := DB.Exec(`create table if not exists authtable(user TEXT, passwd TEXT)`)
	if err != nil {
		logs.ErrorPanic(err, `创建数据表authtable错误`)
	}
	_, err = DB.Exec(`create table if not exists datatable(data BLOB)`)
	if err != nil {
		logs.ErrorPanic(err, `创建数据表datatable错误`)
	}

}
