package base

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//Logger 按天记录日志
type dayLogger struct {
	*log.Logger
}

//LogPath 日志路径
var LogPath = os.Getenv("LOG_PATH")
var (
	_logger    *dayLogger
	loggerOnce sync.Once
)

func init() {
	checkDir(LogPath)
}

//DayLogger 日志
func DayLogger() *dayLogger {
	loggerOnce.Do(func() {
		_logger = &dayLogger{log.New()}
		_logger.Formatter = &log.JSONFormatter{}
	})
	return _logger
}

func checkDir(path string) bool {
	if len(path) == 0 {
		return false
	}

	_, err := os.Stat(path) //os.Stat获取文件信息
	if err == nil {
		return true
	}

	if os.IsExist(err) {
		return true
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.WithFields(log.Fields{"err": err, "path": path}).Error("dir fail")
		return false
	}
	return true
}

func (me *dayLogger) setFileOutput(filename string) bool {
	logFile := fmt.Sprintf("%s/%s-%s.log", LogPath, filename, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "logfile": logFile}).Error("open fail")
		return false
	}
	me.Out = file
	return true
}

func (me *dayLogger) AddRoomLog(action string, data log.Fields, err error) {
	if !me.setFileOutput("liveroom") {
		return
	}

	ctxLogger := me.WithFields(data)
	if err != nil {
		ctxLogger.Error(err)
	} else {
		ctxLogger.Info(action)
	}
}

func (me *dayLogger) AddLog(clientIP, logType string, data interface{}, lerr error) {
	if !me.setFileOutput(logType) {
		return
	}

	ctxLogger := me.WithFields(log.Fields{"ip": clientIP, "type": logType, "data": data})
	if lerr != nil {
		ctxLogger.Error(lerr)
	} else {
		ctxLogger.Info(logType + " success")
	}
}
