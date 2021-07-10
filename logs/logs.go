package logs

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

var errorLog, infoLog = log.New(), log.New()

// errorLog（错误处理）为输出到debug.log，infoLog（消息处理）为输出到terminal
func LogInit() {
	fp, err := os.OpenFile(`error.log`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Panic(errors.New(`[PANIC] Open error.log Failed`))
	}

	// errorLog
	errorLog.Out = fp
	errorLog.Formatter = &log.TextFormatter{
		PadLevelText:           true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	}

	// infoLog
	infoLog.Out = os.Stdout
	infoLog.Formatter = &log.TextFormatter{
		PadLevelText:     true,
		DisableTimestamp: true,
	}
}

// ErrorPanic 处理系统错误，输出到stdout
func ErrorPanic(err error, desc string) {
	infoLog.WithFields(log.Fields{
		`Desc`: desc,
	}).Panicln(err)
	return
}

// ErrorLog 处理一般错误，输出到stdout和log文件
func ErrorLog(err error, desc string) {
	infoLog.WithFields(log.Fields{
		`Desc`: desc,
	}).Errorln(err)
	errorLog.WithFields(log.Fields{
		`Desc`: desc,
	}).Errorln(err)
	return
}
