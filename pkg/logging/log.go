package logging

import (
	"fmt"
	"github.com/op/go-logging"
	"myproject/pmsgo/pkg/setting"
	"os"
)

// 设置日志格式
var logger = logging.MustGetLogger("pms-api")
var format = logging.MustStringFormatter(
	`[%{level:.4s}] %{time:2006-01-02 15:04:05.000} %{message}`,
)

func init() {
	logFile, err := os.OpenFile(setting.LogFile, os.O_WRONLY, 0666)
	if err != nil {
		logFile, err := os.Create(setting.LogFile)
		defer logFile.Close()
		if err != nil {
			fmt.Println(1, err)
		}
	}

	backend1 := logging.NewLogBackend(logFile, "", 0)
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)

	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")
	backend1Formatter := logging.NewBackendFormatter(backend1, format)

	backend2Leveled := logging.AddModuleLevel(backend2)
	backend2Leveled.SetLevel(logging.INFO, "")
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	logging.SetBackend(backend1Formatter, backend2Formatter)
}

func Debug(v ...interface{}) {
	logger.Debug(v)
}

func Info(v ...interface{}) {
	logger.Info(v)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(v ...interface{}) {
	logger.Warning(v)
}

func Error(v ...interface{}) {
	logger.Error(v)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
