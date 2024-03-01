package installer

import (
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
)

var loger logger.Logger = initLogger()

func initLogger() (loger logger.Logger) {
	loger = logrus.Logger
	return
}

func SetLogger(newLoger logger.Logger) {
	loger = newLoger
}
