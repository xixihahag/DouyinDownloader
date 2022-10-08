package logutil

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = log.Sugar()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
