package confs

import (
	"httpserver/logs"
	"os"
	"path"
)

func WdInit() {
	execFile, err := os.Executable()
	if err != nil {
		logs.ErrorPanic(err, `获取运行路径失败`)
	}
	err = os.Chdir(path.Dir(execFile))
	if err != nil {
		logs.ErrorPanic(err, `进入运行路径失败`)
	}
}
