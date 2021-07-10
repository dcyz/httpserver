package confs

import (
	"encoding/json"
	"httpserver/logs"
	"io/ioutil"
)

type NetConf struct {
	Host string
	Port uint32
}

var NetInfo *NetConf

func NetInit() {
	NetInfo = &NetConf{
		Host: `127.0.0.1`,
		Port: 8080,
	}
	NetInfo.Reload()
}

func (n *NetConf) Reload() {
	data, err := ioutil.ReadFile(`./.conf/net.json`)
	if err != nil {
		logs.ErrorPanic(err, `网络配置读取失败`)
	}
	err = json.Unmarshal(data, &NetInfo)
	if err != nil {
		logs.ErrorPanic(err, `网络配置解序列化失败`)
	}
}
