package confs

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kascas/httpserver/logs"
)

type NetConf struct {
	BindIP    string
	PublicIP  string
	PrivateIP string
	Port      uint32
}

var NetInfo *NetConf

func NetInit() {
	NetInfo = &NetConf{
		BindIP:    `127.0.0.1`,
		PublicIP:  `127.0.0.1`,
		PrivateIP: `127.0.0.1`,
		Port:      8080,
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
